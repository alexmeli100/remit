package grpc

import (
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	create            grpcTrans.Handler
	getUserByID       grpcTrans.Handler
	getUserByUUID     grpcTrans.Handler
	getUserByEmail    grpcTrans.Handler
	updateEmail       grpcTrans.Handler
	getContacts       grpcTrans.Handler
	createContact     grpcTrans.Handler
	updateContact     grpcTrans.Handler
	deleteContact     grpcTrans.Handler
	setUserProfile    grpcTrans.Handler
	updateUserProfile grpcTrans.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC NotificatorServer
func NewGRPCServer(endpoints endpoint.Endpoints, options map[string][]grpcTrans.ServerOption) pb.UsersServer {
	return &grpcServer{
		create:            makeCreateHandler(endpoints, options[endpoint.Create]),
		getUserByEmail:    makeGetUserByEmailHandler(endpoints, options[endpoint.GetUserByEmail]),
		getUserByID:       makeGetUserByIDHandler(endpoints, options[endpoint.GetUserById]),
		getUserByUUID:     makeGetUserByUUIDHandler(endpoints, options[endpoint.GetUserByUUID]),
		updateEmail:       makeUpdateEmailHandler(endpoints, options[endpoint.UpdateEmail]),
		createContact:     makeCreateContactHandler(endpoints, options[endpoint.CreateContact]),
		getContacts:       makeGetContactsHandler(endpoints, options[endpoint.GetContacts]),
		deleteContact:     makeDeleteContactHandler(endpoints, options[endpoint.DeleteContact]),
		setUserProfile:    makeSetUserProfileHandler(endpoints, options[endpoint.SetUserProfile]),
		updateUserProfile: makeUpdateUserProfileHandler(endpoints, options[endpoint.UpdateUserProfile]),
		updateContact:     makeUpdateContactHandler(endpoints, options[endpoint.UpdateContact]),
	}
}
