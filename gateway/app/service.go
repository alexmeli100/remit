package app

import (
	"context"
	notificatorService "github.com/alexmeli100/remit/notificator/pkg/service"
	notificatorGRPC "github.com/alexmeli100/remit/notificator/pkg/transport/grpc"
	userGRPC "github.com/alexmeli100/remit/users/pkg/grpc"
	usersService "github.com/alexmeli100/remit/users/pkg/service"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

type GRPCClientOpt = func(map[string][]grpcTrans.ClientOption)

// open connection to user service and return a user service client
func CreateUserServiceClient(ctx context.Context, instance string, options ...GRPCClientOpt) (usersService.UsersService, error) {
	conn, err := grpc.Dial(instance, grpc.WithInsecure())

	go func() {
		defer conn.Close()
		<-ctx.Done()
	}()

	if err != nil {
		return nil, errors.Wrap(err, "error openening connection to user service")
	}

	opts := make(map[string][]grpcTrans.ClientOption)

	for _, option := range options {
		option(opts)
	}

	return userGRPC.NewGRPCClient(conn, opts), nil
}

// open connection to notificator service and return notificator service client
func CreateNotificatorServiceClient(ctx context.Context, instance string, options ...GRPCClientOpt) (notificatorService.NotificatorService, error) {
	conn, err := grpc.Dial(instance, grpc.WithInsecure())

	go func() {
		defer conn.Close()
		<-ctx.Done()
	}()

	if err != nil {
		return nil, errors.Wrap(err, "error openening connection to notificator service")
	}

	opts := make(map[string][]grpcTrans.ClientOption)

	for _, option := range options {
		option(opts)
	}

	return notificatorGRPC.NewGRPCClient(conn, opts), nil
}

// inirialize the app from the options
func (a *App) Initialize(opts ...func(*App) error) error {
	for _, opt := range opts {
		err := opt(a)

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	errc := make(chan error, 1)

	go func() {
		if err := a.Server.ListenAndServe(); err != nil {
			errc <- err
		}
	}()

	select {
	case <-ctx.Done():
		ctx, _ := context.WithTimeout(ctx, time.Second*15)
		return a.Server.Shutdown(ctx)
	case err := <-errc:
		return err
	}
}
