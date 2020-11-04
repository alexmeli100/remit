package main

import (
	"context"
	"fmt"
	"github.com/alexmeli100/remit/events"
	eventpb "github.com/alexmeli100/remit/events/pb"
	"github.com/alexmeli100/remit/payment/cmd/service"
	payService "github.com/alexmeli100/remit/payment/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	"github.com/stripe/stripe-go/v71/client"
	"os"
	"os/signal"
	"time"
)

const (
	NatsClusterId = "nats-streaming-cluster"
	NatsClientId  = "nats-streaming-client"
	DurableName   = "payment-durable"
)

var stripeKey string
var natsInstance string
var logger log.Logger
var dbPassword string
var dbName string
var username string
var dbHost string
var dbPort string

func main() {
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	initFromEnv()

	db, err := openDB(dbHost, dbPort, username, dbPassword, dbName)

	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	if err = db.Ping(); err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	defer func() {
		logger.Log("error", db.Close())
	}()

	sc := &client.API{}
	sc.Init(stripeKey, nil)

	conn, err := events.Connect(natsInstance, NatsClusterId, NatsClientId)

	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	defer func() {
		logger.Log("error", conn.Close())
	}()

	ctx, cancel := context.WithCancel(context.Background())
	postgres := payService.NewPostgresDB(db)
	svc, err := payService.NewStripeService(
		postgres,
		sc,
		logger,
		withUserEventListener(ctx, conn, logger),
		withTransferEventListener(ctx, conn, logger))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		oscall := <-c
		logger.Log("system call", oscall)
		cancel()
	}()

	service.Run(ctx, svc, logger)
}

func openDB(host, port, userName, password, dbName string) (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, userName, password, dbName)

	db, err := sqlx.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func withUserEventListener(ctx context.Context, conn stan.Conn, logger log.Logger) func(*payService.StripeService) error {
	return func(stripeService *payService.StripeService) error {
		handlers := map[eventpb.EventKind]events.EventHandler{
			eventpb.UserCreated: stripeService.OnUserCreated,
		}

		err := addEventListener(ctx, conn, logger, handlers, events.UserEvents, "payment-user-events")

		return err
	}
}

func withTransferEventListener(ctx context.Context, conn stan.Conn, logger log.Logger) func(*payService.StripeService) error {
	return func(stripeService *payService.StripeService) error {
		handlers := map[eventpb.EventKind]events.EventHandler{
			eventpb.TransferSucceded: stripeService.OnTransferSucceded,
		}

		err := addEventListener(ctx, conn, logger, handlers, events.TransferEvents, "payment-transfer-events")

		return err
	}
}

func addEventListener(ctx context.Context, conn stan.Conn, logger log.Logger, handlers events.Handlers, topic, queue string) error {
	errc, err := events.ListenEvents(
		ctx,
		topic,
		queue,
		conn,
		handlers,
		stan.SetManualAckMode(),
		stan.AckWait(time.Minute),
		stan.DurableName(DurableName),
		stan.MaxInflight(25))

	if err != nil {
		return err
	}

	go func() {
		for err = range errc {
			logger.Log("topic", topic, "err", err)
		}
	}()

	return nil
}

func initFromEnv() {
	stripeKey = os.Getenv("STRIPE_API_KEY")
	natsHost := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")
	natsInstance = natsHost + ":" + natsPort
	dbPassword = os.Getenv("POSTGRES_PASSWORD")
	username = os.Getenv("POSTGRES_USER")
	dbName = os.Getenv("POSTGRES_DB")
	dbHost = os.Getenv("PAYMENT_DB_SERVICE_HOST")
	dbPort = os.Getenv("PAYMENT_DB_SERVICE_PORT")
}
