package grpc

import (
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	create         grpcTrans.Handler
	getUserByID    grpcTrans.Handler
	getUserByUUID  grpcTrans.Handler
	getUserByEmail grpcTrans.Handler
	updateEmail    grpcTrans.Handler
	updateStatus   grpcTrans.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC NotificatorServer
func NewGRPCServer(endpoints endpoint.Endpoints, options map[string][]grpcTrans.ServerOption) pb.UsersServer {
	return &grpcServer{
		create:         makeCreateHandler(endpoints, options[endpoint.Create]),
		getUserByEmail: makeGetUserByEmailHandler(endpoints, options[endpoint.GetUserByEmail]),
		getUserByID:    makeGetUserByIDHandler(endpoints, options[endpoint.GetUserById]),
		getUserByUUID:  makeGetUserByUUIDHandler(endpoints, options[endpoint.GetUserByUUID]),
		updateEmail:    makeUpdateEmailHandler(endpoints, options[endpoint.UpdateEmail]),
		updateStatus:   makeUpdateStatusHandler(endpoints, options[endpoint.UpdateStatus]),
	}
}
