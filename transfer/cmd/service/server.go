package service

import (
	"context"
	endpoint "github.com/alexmeli100/remit/transfer/pkg/endpoint"
	transferTran "github.com/alexmeli100/remit/transfer/pkg/grpc"
	"github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	log "github.com/go-kit/kit/log"
	opentracing "github.com/go-kit/kit/tracing/opentracing"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type TransferServerOpt = map[string][]grpcTrans.ServerOption
type TransferServerOptFunc = func(opt TransferServerOpt)

func createService(endpoints endpoint.Endpoints, opts ...TransferServerOptFunc) *grpc.Server {
	serverOptions := make(TransferServerOpt)

	for _, f := range opts {
		f(serverOptions)
	}

	grpcServer := transferTran.NewGRPCServer(endpoints, serverOptions)
	baseServer := grpc.NewServer()
	pb.RegisterTransferServer(baseServer, grpcServer)
	reflection.Register(baseServer)

	return baseServer
}

func runServer(ctx context.Context, server *grpc.Server) error {
	grpcListener, err := net.Listen("tcp", grpcAddr)
	errc := make(chan error, 1)

	if err != nil {
		return err
	}

	go func() {
		if err = server.Serve(grpcListener); err != nil {
			errc <- err
		}
	}()

	select {
	case <-ctx.Done():
		server.GracefulStop()
	case err := <-errc:
		return errors.Wrap(err, "listen error")
	}

	return nil
}

// add a logger to the server endpoints
func withLogger(logger log.Logger, endpoints ...string) TransferServerOptFunc {
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
func withTracer(tracer stdopentracing.Tracer, logger log.Logger, endpoints ...string) TransferServerOptFunc {
	return func(m TransferServerOpt) {
		for _, ep := range endpoints {
			e, ok := m[ep]

			if !ok {
				e = make([]grpcTrans.ServerOption, 0, 2)
			}

			m[ep] = append(e, grpcTrans.ServerBefore(opentracing.GRPCToContext(tracer, ep, logger)))
		}
	}
}
