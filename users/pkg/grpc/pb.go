package grpc

import (
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/alexmeli100/remit/users/pkg/service"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer
type grpcServer struct {
	create         grpcTrans.Handler
	getUserByID    grpcTrans.Handler
	getUserByUUID  grpcTrans.Handler
	getUserByEmail grpcTrans.Handler
	updateEmail    grpcTrans.Handler
	updateStatus   grpcTrans.Handler
}

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

func NewGRPCClient(conn *grpc.ClientConn, options map[string][]grpcTrans.ClientOption) service.UsersService {
	return endpoint.Endpoints{
		CreateEndpoint:         makeCreateClient(conn, options[endpoint.Create]).Endpoint(),
		GetUserByEmailEndpoint: makeGetUserByEmailClient(conn, options[endpoint.GetUserByEmail]).Endpoint(),
		GetUserByIDEndpoint:    makeGetUserByIDClient(conn, options[endpoint.GetUserById]).Endpoint(),
		GetUserByUUIDEndpoint:  makeGetUserByUUIDClient(conn, options[endpoint.GetUserByUUID]).Endpoint(),
		UpdateEmailEndpoint:    makeUpdateEmailClient(conn, options[endpoint.UpdateEmail]).Endpoint(),
		UpdateStatusEndpoint:   makeUpdateStatusClient(conn, options[endpoint.UpdateStatus]).Endpoint(),
	}
}
