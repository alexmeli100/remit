package service

import (
	"context"
	endpoint "github.com/alexmeli100/remit/payment/pkg/endpoint"
	service "github.com/alexmeli100/remit/payment/pkg/service"
	kitEndpoint "github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	opentracinggo "github.com/opentracing/opentracing-go"
	"os"
)

var tracer opentracinggo.Tracer
var grpcAddr string

func Run(ctx context.Context, p service.PaymentService, logger log.Logger) {
	port := os.Getenv("Payment_MANAGER_SERVICE_PORT")
	grpcAddr = ":" + port
	tracer = opentracinggo.GlobalTracer()

	el := endpoint.GetEndpointList()
	svc := getService(p, serviceWithLogger(logger))
	eps := getEndpoint(svc, endpointWithLogger(logger, el...))
	server := createService(eps, withLogger(logger, el...), withTracer(tracer, logger, el...))

	logger.Log("transport", "gRPC", "addr", grpcAddr)
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
func getService(svc service.PaymentService, opts ...func([]service.Middleware) []service.Middleware) service.PaymentService {
	mw := make([]service.Middleware, 0, 4)

	for _, opt := range opts {
		mw = opt(mw)
	}

	return service.New(svc, mw)
}

// add the middlewares and get the endpoints from the payment service
func getEndpoint(svc service.PaymentService, opts ...func(map[string][]kitEndpoint.Middleware)) endpoint.Endpoints {
	mw := make(map[string][]kitEndpoint.Middleware)

	for _, opt := range opts {
		opt(mw)
	}

	return endpoint.New(svc, mw)
}
