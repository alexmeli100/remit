package service

import (
	"flag"
	"fmt"
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/alexmeli100/remit/users/pkg/grpc"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/alexmeli100/remit/users/pkg/service"
	"net"
	"os"
	"os/signal"
	"syscall"

	endpoint1 "github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	group "github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
	grpc1 "google.golang.org/grpc"
)

var tracer opentracinggo.Tracer
var logger log.Logger

var fs = flag.NewFlagSet("users", flag.ExitOnError)
var grpcAddr = fs.String("grpc-addr", ":8082", "gRPC listen address")

func Run() {
	fs.Parse(os.Args[1:])

	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	tracer = opentracinggo.GlobalTracer()

	svc := service.New(getServiceMiddleware(logger))
	eps := endpoint.New(svc, getEndpointMiddleware(logger))
	g := createService(eps)
	initCancelInterrupt(g)
	logger.Log("exit", g.Run())
}

func getServiceMiddleware(logger log.Logger) []service.Middleware {
	var mw []service.Middleware
	mw = append(mw, service.LoggingMiddleware(logger))

	return mw
}
func getEndpointMiddleware(logger log.Logger) map[string][]endpoint1.Middleware {
	mw := make(map[string][]endpoint1.Middleware)
	endpoint.AddDefaultEndPointMiddleware(logger, mw)

	return mw
}

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})

	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

func initGRPCHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultGRPCOptions(logger, tracer)

	grpcServer := grpc.NewGRPCServer(endpoints, options)
	grpcListener, err := net.Listen("tcp", *grpcAddr)

	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	}

	g.Add(func() error {
		logger.Log("transport", "gRPC", "addr", *grpcAddr)
		baseServer := grpc1.NewServer()
		pb.RegisterUsersServer(baseServer, grpcServer)

		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})
}
