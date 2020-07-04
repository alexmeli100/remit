package service

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	userTran "github.com/alexmeli100/remit/users/pkg/grpc"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

type UserServerOpt = map[string][]grpcTrans.ServerOption
type UserServerOptFunc = func(opt UserServerOpt)

// create a notificator server from the endpoints and apply the server options
func createService(endpoints endpoint.Endpoints, opts ...UserServerOptFunc) *grpc.Server {
	serverOptions := make(UserServerOpt)

	for _, f := range opts {
		f(serverOptions)
	}

	grpcServer := userTran.NewGRPCServer(endpoints, serverOptions)
	baseServer := grpc.NewServer()
	pb.RegisterUsersServer(baseServer, grpcServer)

	return baseServer
}

func runServer(ctx context.Context, server *grpc.Server) error {
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	errc := make(chan error, 1)

	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		return err
	}

	go func() {
		if err = server.Serve(grpcListener); err != nil {
			errc <- err
		}
	}()

	logger.Log("transport", "gRPC", "addr", *grpcAddr)

	select {
	case <-ctx.Done():
		server.GracefulStop()
	case err := <-errc:
		return errors.Wrap(err, "listen error")
	}

	return nil
}

// add a logger to the server endpoints
func withLogger(logger log.Logger, endpoints ...string) UserServerOptFunc {
	return func(m map[string][]grpcTrans.ServerOption) {
		for _, ep := range endpoints {
			e, ok := m[ep]

			if !ok {
				e = make([]grpcTrans.ServerOption, 0, 2)
			}

			m[ep] = append(e, grpcTrans.ServerErrorLogger(logger))
		}
	}
}

// add a tracer to the server endpoints
func withTracer(tracer stdopentracing.Tracer, logger log.Logger, endpoints ...string) UserServerOptFunc {
	return func(m UserServerOpt) {
		for _, ep := range endpoints {
			e, ok := m[ep]

			if !ok {
				e = make([]grpcTrans.ServerOption, 0, 2)
			}

			m[ep] = append(e, grpcTrans.ServerBefore(opentracing.GRPCToContext(tracer, ep, logger)))
		}
	}
}
