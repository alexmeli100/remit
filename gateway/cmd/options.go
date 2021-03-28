package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/alexmeli100/remit/events"
	"github.com/alexmeli100/remit/gateway/app"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/cors"
	"google.golang.org/api/option"
	"net/http"
	"time"
)

func svcWithTracer(tracer stdopentracing.Tracer, logger log.Logger, endpoints ...string) app.GRPCClientOpt {
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
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"X-Requested-With", "Authorization", "Content-Type"},
			AllowCredentials: true,
		})

		s.Handler = c.Handler(r)
	}
}

func serverWithIdleTimeout(t time.Duration) func(server *http.Server) {
	return func(s *http.Server) {
		s.IdleTimeout = t
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
func appWithPaymentService(ctx context.Context, instance string, opts ...app.GRPCClientOpt) func(*app.App) error {
	return func(a *app.App) error {
		p, err := app.CreatePaymentServiceClient(ctx, instance, opts...)

		if err != nil {
			return err
		}

		a.PaymentService = p
		return nil
	}
}

func appWithTransferService(ctx context.Context, instance string, opts ...app.GRPCClientOpt) func(*app.App) error {
	return func(a *app.App) error {
		t, err := app.CreateTransferServiceClient(ctx, instance, opts...)

		if err != nil {
			return err
		}

		a.TransferService = t
		return nil
	}
}

// add a listener for user events
func appWithUserEventListener(ctx context.Context, conn stan.Conn, logger log.Logger) func(*app.App) error {
	return func(a *app.App) error {
		errc, err := events.ListenAllUserEvents(
			ctx,
			conn,
			"user-events-gateway",
			a,
			stan.SetManualAckMode(),
			stan.AckWait(time.Minute),
			stan.DurableName(DurableName),
			stan.MaxInflight(25))

		if err != nil {
			return err
		}

		// log errors from user event listener
		go func() {
			for err = range errc {
				logger.Log("method", "listen-user-events", "err", err)
			}
		}()

		return nil
	}
}

func appWithEventSender(ctx context.Context, conn stan.Conn) func(*app.App) error {
	return func(a *app.App) error {
		sender := events.NewEventSender(conn)

		a.Events = sender
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

func appWithRedisClient(ctx context.Context, options *redis.Options) func(*app.App) error {
	return func(a *app.App) error {
		client := redis.NewClient(options)

		_, err := client.Ping(ctx).Result()

		if err != nil {
			return err
		}

		go func() {
			<-ctx.Done()
			a.Logger.Log("error", client.Close())
		}()

		a.RedisClient = client
		return nil
	}
}

func appWithLogger(logger log.Logger) func(*app.App) error {
	return func(a *app.App) error {
		a.Logger = logger
		return nil
	}
}
