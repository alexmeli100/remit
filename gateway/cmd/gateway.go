package main

import (
	"context"
	"flag"
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
var fs = flag.NewFlagSet("gateway", flag.ExitOnError)

func main() {
	fs.Parse(os.Args[1:])

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
		serverWithAddress(":8085"),
		serverWithHandler(r),
		serverWithReadTimeout(time.Second*15),
		serverWithWriteTimeout(time.Second*15))

	err := a.Initialize(
		serverFunc,
		appWithFirebase(ctx, "firebase-service-account.json"),
		appWithUserService(ctx, ":8081", userWithTracer(tracer, logger, el...)),
		appWithNotificatorService(ctx, ":8083", notificatorWithTracer(tracer, logger, el...)))

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

	logger.Log("exit", a.Run(ctx))
}
