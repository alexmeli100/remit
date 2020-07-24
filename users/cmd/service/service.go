package service

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/alexmeli100/remit/users/pkg/service"
	"os"
	"os/signal"

	kitEndpoint "github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	opentracinggo "github.com/opentracing/opentracing-go"
)

var tracer opentracinggo.Tracer
var logger log.Logger
var grpcAddr string

func Run(u service.UsersService) {
	port := os.Getenv("USER_MANAGER_SERVICE_PORT")

	grpcAddr = ":" + port

	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	tracer = opentracinggo.GlobalTracer()
	logger.Log("grpcAddr", grpcAddr)

	el := endpoint.GetEndpointList()
	svc := getService(u, serviceWithLogger(logger))
	eps := getEndpoint(svc, endpointWithLogger(logger, el...))
	server := createService(eps, withLogger(logger, el...), withTracer(tracer, logger, el...))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		logger.Log("system call", oscall)
		cancel()
	}()

	logger.Log("exit", runServer(ctx, server))
}

// add logger to notificator service
func serviceWithLogger(logger log.Logger) func([]service.Middleware) []service.Middleware {
	return func(mw []service.Middleware) []service.Middleware {
		mw = append(mw, service.LoggingMiddleware(logger))

		return mw
	}
}

// add logger to endpoint
func endpointWithLogger(logger log.Logger, eps ...string) func(map[string][]kitEndpoint.Middleware) {
	return func(mw map[string][]kitEndpoint.Middleware) {
		for _, name := range eps {
			logMw := endpoint.LoggingMiddleware(logger)

			m, ok := mw[name]

			if !ok {
				m = make([]kitEndpoint.Middleware, 0, 2)
			}

			mw[name] = append(m, logMw)
		}
	}
}

// add the middlewares and get the user service
func getService(svc service.UsersService, opts ...func([]service.Middleware) []service.Middleware) service.UsersService {
	mw := make([]service.Middleware, 0, 4)

	for _, opt := range opts {
		mw = opt(mw)
	}

	return service.New(svc, mw)
}

// add the middlewares and get the endpoints from the user service
func getEndpoint(svc service.UsersService, opts ...func(map[string][]kitEndpoint.Middleware)) endpoint.Endpoints {
	mw := make(map[string][]kitEndpoint.Middleware)

	for _, opt := range opts {
		opt(mw)
	}

	return endpoint.New(svc, mw)
}
