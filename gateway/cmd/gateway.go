package main

import (
	"context"
	"fmt"
	"github.com/alexmeli100/remit/gateway/app"
	notificator "github.com/alexmeli100/remit/notificator/pkg/service"
	transfer "github.com/alexmeli100/remit/transfer/pkg/service"
	user "github.com/alexmeli100/remit/users/pkg/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	opentracinggo "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"time"
)

var tracer opentracinggo.Tracer

func openDB(host, port, userName, password, dbName string) (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, userName, password, dbName)

	db, err := sqlx.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	// Configure logger
	logger, err := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	tracer = opentracinggo.GlobalTracer()

	ctx, cancel := context.WithCancel(context.Background())
	a := &app.App{}
	r := mux.NewRouter()
	a.InitializeRoutes(r)

	password := os.Getenv("POSTGRES_PASSWORD")
	userName := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("USER_DB_SERVICE_HOST")
	dbPort := os.Getenv("USER_DB_SERVICE_PORT")

	db, err := openDB(dbHost, dbPort, userName, password, dbName)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	serverFunc := appWithServer(
		serverWithAddress(":8083"),
		serverWithHandler(r),
		serverWithReadTimeout(time.Second*5),
		serverWithIdleTimeout(time.Second*10),
		serverWithWriteTimeout(time.Second*15))

	// Get the user service
	userSVC := appWithUserPostgService(ctx, db, user.LoggingMiddleware(logger))

	// Get the notification service
	sendGridApiKey := os.Getenv("SENDGRID_API_KEY")
	notificatorSVC := appWithNotificatorService(
		ctx,
		sendGridApiKey,
		notificator.LoggingMiddleware(logger))

	// Get the transfer service
	transferSVC := appWithTransferService(ctx, transfer.LoggingMiddleware(logger))

	err = a.Initialize(
		serverFunc,
		userSVC,
		transferSVC,
		notificatorSVC,
		appWithFirebase(ctx, "/opt/firebase/wealow-test-firebase.json"),
		appWithLogger(logger))

	if err != nil {
		logger.Error("Error initiazing app", zap.Error(err))
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		oscall := <-c
		logger.Info("system call", zap.String("signal", oscall.String()))
		cancel()
	}()

	logger.Info("addr", zap.String("Server addr", a.Server.Addr))
	logger.Info("exit", zap.String("exit status", zap.Error(a.Run(ctx)).String))
}
