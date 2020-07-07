package grpc

import (
	"github.com/alexmeli100/remit/notificator/pkg/endpoint"
	"github.com/alexmeli100/remit/notificator/pkg/service"
	"github.com/alexmeli100/remit/notificator/pkg/transport/grpc/pb"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type grpcServer struct {
	sendConfirmEmail       grpcTrans.Handler
	sendPasswordResetEmail grpcTrans.Handler
	sendWelcomeEmail       grpcTrans.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC NotificatorServer
func NewGRPCServer(endpoints endpoint.Endpoints, options map[string][]grpcTrans.ServerOption) pb.NotificatorServer {
	return &grpcServer{
		sendConfirmEmail:       makeSendConfirmEmailHandler(endpoints, options[endpoint.SendConfirmEmail]),
		sendPasswordResetEmail: makeSendPasswordResetEmailHandler(endpoints, options[endpoint.SendPasswordResetEmail]),
		sendWelcomeEmail:       makeSendWelcomeEmailHandler(endpoints, options[endpoint.SendWelcomeEmail]),
	}
}

func NewGRPCClient(conn *grpc.ClientConn, options map[string][]grpcTrans.ClientOption) service.NotificatorService {
	return endpoint.Endpoints{
		SendConfirmEmailEndpoint:       makeSendConfirmEmailClient(conn, options[endpoint.SendConfirmEmail]).Endpoint(),
		SendPasswordResetEmailEndpoint: makeSendPasswordResetEmailClient(conn, options[endpoint.SendPasswordResetEmail]).Endpoint(),
		SendWelcomeEmailEndpoint:       makeSendWelcomeEmailClient(conn, options[endpoint.SendWelcomeEmail]).Endpoint(),
	}
}
