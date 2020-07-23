package main

import (
	"context"
	"github.com/alexmeli100/remit/gateway/app"
	userEndpoint "github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	opentracinggo "github.com/opentracing/opentracing-go"
	"os"
	"os/signal"
	"time"
)

var tracer opentracinggo.Tracer
var logger log.Logger
var usersInstance string
var notificatorInstance string
var natsInstance string

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
	el := userEndpoint.GetEndpointList()

	serverFunc := appWithServer(
		serverWithAddress(":5000"),
		serverWithHandler(r),
		serverWithReadTimeout(time.Second*15),
		serverWithWriteTimeout(time.Second*15))

	err := a.Initialize(
		serverFunc,
		appWithFirebase(ctx, "/opt/firebase/firebase-service-account.json"),
		appWithEventSender(ctx, natsInstance),
		appWithUserEventListener(ctx, natsInstance),
		appWithUserService(ctx, usersInstance, userWithTracer(tracer, logger, el...)),
		appWithNotificatorService(ctx, notificatorInstance, notificatorWithTracer(tracer, logger, el...)))

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

	logger.Log("server running at", a.Server.Addr)
	logger.Log("exit", a.Run(ctx))
}

func initFromEnv() {
	usersHost := os.Getenv("USER_MANAGER_SERVICE_HOST")
	usersPort := os.Getenv("USER_MANAGER_SERVICE_PORT")
	notificatorHost := os.Getenv("NOTIFICATOR_MANAGER_SERVICE_HOST")
	notificatorPort := os.Getenv("NOTIFICATOR_MANAGER_SERVICE_PORT")
	natsHost := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")
	usersInstance = usersHost + ":" + usersPort
	notificatorInstance = notificatorHost + ":" + notificatorPort
	natsInstance = natsHost + ":" + natsPort
}
