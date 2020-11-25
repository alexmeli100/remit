package main

import (
	"context"
	"github.com/alexmeli100/remit/events"
	"github.com/alexmeli100/remit/gateway/app"
	notificatorEndpoint "github.com/alexmeli100/remit/notificator/pkg/endpoint"
	paymentEndpoint "github.com/alexmeli100/remit/payment/pkg/endpoint"
	userEndpoint "github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	opentracinggo "github.com/opentracing/opentracing-go"
	"os"
	"os/signal"
	"time"
)

const (
	NatsClusterId = "nats-streaming-cluster"
	NatsClientId  = "nats-streaming-client-2"
	DurableName   = "gateway-durable"
)

var tracer opentracinggo.Tracer
var logger log.Logger
var usersInstance string
var notificatorInstance string
var paymentInstance string
var natsInstance string
var redisInstance string
var redisPass string

func main() {
	initFromEnv()

	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	tracer = opentracinggo.GlobalTracer()

	ctx, cancel := context.WithCancel(context.Background())
	a := &app.App{}
	r := mux.NewRouter()
	a.InitializeRoutes(r)
	userEl := userEndpoint.GetEndpointList()
	paymentEl := paymentEndpoint.GetEndpointList()
	notificatorEl := notificatorEndpoint.GetEndpointList()

	conn, err := events.Connect(natsInstance, NatsClusterId, NatsClientId)

	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	defer func() {
		logger.Log("error", conn.Close())
	}()

	redisOptions := &redis.Options{
		Addr:     redisInstance,
		Password: redisPass,
	}

	serverFunc := appWithServer(
		serverWithAddress(":8083"),
		serverWithHandler(r),
		serverWithReadTimeout(time.Second*5),
		serverWithIdleTimeout(time.Second*10),
		serverWithWriteTimeout(time.Second*15))

	userSVC := appWithUserService(ctx,
		usersInstance,
		svcWithTracer(tracer, logger, userEl...))

	paymentSVC := appWithPaymentService(
		ctx,
		paymentInstance,
		svcWithTracer(tracer, logger, paymentEl...))

	notificatorSVC := appWithNotificatorService(
		ctx,
		notificatorInstance,
		svcWithTracer(tracer, logger, notificatorEl...))

	err = a.Initialize(
		serverFunc,
		userSVC,
		paymentSVC,
		notificatorSVC,
		appWithFirebase(ctx, "/opt/firebase/wealow-test-firebase.json"),
		appWithEventSender(ctx, conn),
		appWithLogger(logger),
		appWithRedisClient(ctx, redisOptions),
		appWithUserEventListener(ctx, conn, logger))

	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		oscall := <-c
		logger.Log("system call", oscall)
		cancel()
	}()

	logger.Log("addr", a.Server.Addr)
	logger.Log("exit", a.Run(ctx))
}

func initFromEnv() {
	usersHost := os.Getenv("USER_MANAGER_SERVICE_HOST")
	usersPort := os.Getenv("USER_MANAGER_SERVICE_PORT")
	notificatorHost := os.Getenv("NOTIFICATOR_MANAGER_SERVICE_HOST")
	notificatorPort := os.Getenv("NOTIFICATOR_MANAGER_SERVICE_PORT")
	paymentHost := os.Getenv("PAYMENT_MANAGER_SERVICE_HOST")
	paymentPort := os.Getenv("PAYMENT_MANAGER_SERVICE_PORT")
	natsHost := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPass = os.Getenv("REDIS_PASSWORD")
	usersInstance = usersHost + ":" + usersPort
	notificatorInstance = notificatorHost + ":" + notificatorPort
	natsInstance = natsHost + ":" + natsPort
	paymentInstance = paymentHost + ":" + paymentPort
	redisInstance = redisHost + ":" + redisPort
}
