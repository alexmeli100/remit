package grpc

import (
	"github.com/alexmeli100/remit/notificator/pkg/endpoint"
	"github.com/alexmeli100/remit/notificator/pkg/transport/grpc/pb"
	grpc "github.com/go-kit/kit/transport/grpc"
)

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer
type grpcServer struct {
	sendConfirmEmail       grpc.Handler
	sendPasswordResetEmail grpc.Handler
	sendWelcomeEmail       grpc.Handler
}

func NewGRPCServer(endpoints endpoint.Endpoints, options map[string][]grpc.ServerOption) pb.NotificatorServer {
	return &grpcServer{
		sendConfirmEmail:       makeSendConfirmEmailHandler(endpoints, options[endpoint.SendConfirmEmail]),
		sendPasswordResetEmail: makeSendPasswordResetEmailHandler(endpoints, options[endpoint.SendPasswordResetEmail]),
		sendWelcomeEmail:       makeSendWelcomeEmailHandler(endpoints, options[endpoint.SendWelcomeEmail]),
	}
}
