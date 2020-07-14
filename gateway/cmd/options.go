package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/alexmeli100/remit/gateway/app"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/api/option"
	"net/http"
	"time"
)

// user service client options
func userWithTracer(tracer stdopentracing.Tracer, logger log.Logger, endpoints ...string) app.GRPCClientOpt {
	return func(m map[string][]grpcTrans.ClientOption) {
		for _, endpoint := range endpoints {
			e, ok := m[endpoint]

			if !ok {
				e = make([]grpcTrans.ClientOption, 0)
			}

			m[endpoint] = append(e, grpcTrans.ClientBefore(opentracing.ContextToGRPC(tracer, logger)))
		}
	}
}

// notificator service client options
func notificatorWithTracer(tracer stdopentracing.Tracer, logger log.Logger, endpoints ...string) app.GRPCClientOpt {
	return func(m map[string][]grpcTrans.ClientOption) {
		for _, endpoint := range endpoints {
			e, ok := m[endpoint]

			if !ok {
				e = make([]grpcTrans.ClientOption, 0)
			}

			m[endpoint] = append(e, grpcTrans.ClientBefore(opentracing.ContextToGRPC(tracer, logger)))
		}
	}
}

// server options
func serverWithWriteTimeout(t time.Duration) func(*http.Server) {
	return func(s *http.Server) {
		s.WriteTimeout = t
	}
}

func serverWithReadTimeout(t time.Duration) func(*http.Server) {
	return func(s *http.Server) {
		s.ReadTimeout = t
	}
}

func serverWithAddress(addr string) func(*http.Server) {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

func serverWithHandler(r *mux.Router) func(*http.Server) {
	return func(s *http.Server) {
		s.Handler = r
	}
}

// application options
func appWithServer(opts ...func(*http.Server)) func(*app.App) error {
	return func(app *app.App) error {
		s := &http.Server{}

		for _, opt := range opts {
			opt(s)
		}

		app.Server = s
		return nil
	}
}

func appWithUserService(ctx context.Context, instance string, opts ...app.GRPCClientOpt) func(*app.App) error {
	return func(a *app.App) error {
		u, err := app.CreateUserServiceClient(ctx, instance, opts...)

		if err != nil {
			return err
		}

		a.UsersService = u
		return nil
	}
}

func appWithNotificatorService(ctx context.Context, instance string, opts ...app.GRPCClientOpt) func(*app.App) error {
	return func(a *app.App) error {
		n, err := app.CreateNotificatorServiceClient(ctx, instance, opts...)

		if err != nil {
			return err
		}

		a.Notificator = n
		return nil
	}
}

func appWithFirebase(ctx context.Context, service string) func(*app.App) error {
	return func(a *app.App) error {
		opt := option.WithCredentialsFile(service)

		fireApp, err := firebase.NewApp(ctx, nil, opt)

		if err != nil {
			return err
		}

		a.FireApp = fireApp
		return nil
	}
}