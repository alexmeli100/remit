package service

import (
	"context"
	notificatorEndpoint "github.com/alexmeli100/remit/notificator/pkg/endpoint"
	notificatorTrans "github.com/alexmeli100/remit/notificator/pkg/transport/grpc"
	"github.com/alexmeli100/remit/notificator/pkg/transport/grpc/pb"
	log "github.com/go-kit/kit/log"
	opentracing "github.com/go-kit/kit/tracing/opentracing"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type NotificatorServerOpt = map[string][]grpcTrans.ServerOption
type NotificatorServerOptFunc = func(opt NotificatorServerOpt)

// create a notificator server from the endpoints and apply the server options
func createService(endpoints notificatorEndpoint.Endpoints, opts ...NotificatorServerOptFunc) *grpc.Server {
	serverOptions := make(NotificatorServerOpt)

	for _, f := range opts {
		f(serverOptions)
	}

	grpcServer := notificatorTrans.NewGRPCServer(endpoints, serverOptions)
	baseServer := grpc.NewServer()
	pb.RegisterNotificatorServer(baseServer, grpcServer)

	return baseServer
}

// run the server until an error is thrown or the context gets cancelled
func runServer(ctx context.Context, server *grpc.Server) error {
	grpcListener, err := net.Listen("tcp", grpcAddr)
	errc := make(chan error, 1)
	reflection.Register(server)

	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		return err
	}

	go func() {
		if err = server.Serve(grpcListener); err != nil {
			errc <- err
		}
	}()

	logger.Log("transport", "gRPC", "addr", grpcAddr)

	select {
	case <-ctx.Done():
		server.GracefulStop()
	case err := <-errc:
		return errors.Wrap(err, "listen error")
	}

	return nil
}

// add a logger to the server endpoints
func withLogger(logger log.Logger, endpoints ...string) NotificatorServerOptFunc {
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
func withTracer(tracer stdopentracing.Tracer, logger log.Logger, endpoints ...string) NotificatorServerOptFunc {
	return func(m NotificatorServerOpt) {
		for _, ep := range endpoints {
			e, ok := m[ep]

			if !ok {
				e = make([]grpcTrans.ServerOption, 0, 2)
			}

			m[ep] = append(e, grpcTrans.ServerBefore(opentracing.GRPCToContext(tracer, ep, logger)))
		}
	}
}
