package client

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/alexmeli100/remit/users/pkg/service"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func makeCreateClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "Create", encodeCreateRequest, decodeCreateResponse, pb.CreateReply{}, options...)
}

func encodeCreateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.CreateRequest)

	return &pb.CreateRequest{User: req.User}, nil
}

func decodeCreateResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.CreateReply)

	return endpoint.CreateResponse{Err: str2err(res.Err)}, nil
}

func makeGetUserByIDClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "GetUserByID", encodeGetUserByIDRequest, decodeGetUserByIDResponse, pb.GetUserByIDReply{}, options...)
}

func decodeGetUserByIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetUserByIDReply)

	return endpoint.GetUserByIDResponse{User: res.User, Err: str2err(res.Err)}, nil
}

func encodeGetUserByIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.GetUserByIDRequest)

	return &pb.GetUserByIDRequest{Id: req.Id}, nil
}

func makeGetUserByUUIDClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "GetUserByUUID", encodeGetUserByUUIDRequest, decodeGetUserByUUIDResponse, pb.GetUserByUUIDReply{}, options...)
}

func decodeGetUserByUUIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetUserByUUIDReply)

	return endpoint.GetUserByUUIDResponse{User: res.User, Err: str2err(res.Err)}, nil
}

func encodeGetUserByUUIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.GetUserByUUIDRequest)

	return &pb.GetUserByUUIDRequest{UUID: req.UUID}, nil
}

func makeGetUserByEmailClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "GetUserByEmail", encodeGetUserByEmailRequest, decodeGetUserByEmailResponse, pb.GetUserByEmailReply{}, options...)
}

func decodeGetUserByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetUserByEmailReply)

	return endpoint.GetUserByEmailResponse{User: res.User, Err: str2err(res.Err)}, nil
}

func encodeGetUserByEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.GetUserByEmailRequest)

	return &pb.GetUserByEmailRequest{Email: req.Email}, nil
}

func makeUpdateEmailClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Users",
		"UpdateEmail",
		encodeUpdateEmailRequest,
		decodeUpdateEmailResponse,
		pb.UpdateEmailReply{},
		options...)
}

func decodeUpdateEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.UpdateEmailReply)

	return endpoint.UpdateEmailResponse{Err: str2err(res.Err)}, nil
}

func encodeUpdateEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.UpdateEmailRequest)

	return &pb.UpdateEmailRequest{User: req.User}, nil
}

func makeUpdateContactClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Users",
		"UpdateContact",
		encodeUpdateContactRequest,
		decodeUpdateContactResponse,
		pb.UpdateContactReply{},
		options...)
}

func encodeUpdateContactRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.UpdateContactRequest)

	return &pb.UpdateContactRequest{Contact: req.Contact}, nil
}

func decodeUpdateContactResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.UpdateContactReply)

	return endpoint.UpdateContactResponse{Contact: res.Contact, Err: str2err(res.Err)}, nil
}

func makeUpdateUserProfileClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Users",
		"UpdateUserProfile",
		encodeUpdateUserProfileRequest,
		decodeUpdateUserProfileResponse,
		pb.UpdateUserProfileReply{},
		options...)
}

func encodeUpdateUserProfileRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.UpdateUserProfileRequest)

	return &pb.UpdateUserProfileRequest{User: req.User}, nil
}

func decodeUpdateUserProfileResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.UpdateUserProfileReply)

	return endpoint.UpdateUserProfileResponse{User: res.User, Err: str2err(res.Err)}, nil
}

func makeSetUserProfileClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Users",
		"SetUserProfile",
		encodeSetUserProfileRequest,
		decodeSetUserProfileResponse,
		pb.SetUserProfileReply{},
		options...)
}

func encodeSetUserProfileRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.SetUserProfileRequest)

	return &pb.SetUserProfileRequest{User: req.User}, nil
}

func decodeSetUserProfileResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.SetUserProfileReply)

	return endpoint.SetUserProfileResponse{User: res.User, Err: str2err(res.Err)}, nil
}

func makeDeleteContactClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Users",
		"DeleteContact",
		encodeDeleteContactRequest,
		decodeDeleteContactResponse,
		pb.DeleteContactReply{},
		options...)
}

func encodeDeleteContactRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.DeleteContactRequest)

	return &pb.DeleteContactRequest{Contact: req.Contact}, nil
}

func decodeDeleteContactResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.DeleteContactReply)

	return endpoint.DeleteContactResponse{Err: str2err(res.Err)}, nil
}

func makeCreateContactClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Users",
		"CreateContact",
		encodeCreateContactRequest,
		decodeCreateContactResponse,
		pb.CreateContactReply{},
		options...)
}

func encodeCreateContactRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.CreateContactRequest)

	return &pb.CreateContactRequest{Contact: req.Contact}, nil
}

func decodeCreateContactResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.CreateContactReply)

	return endpoint.CreateContactResponse{Contact: res.Contact, Err: str2err(res.Err)}, nil
}

func makeGetContactsClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Users",
		"GetContacts",
		encodeGetContactsRequest,
		decodeGetContactsResponse,
		pb.GetContactsReply{},
		options...)
}

func encodeGetContactsRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.GetContactsRequest)

	return &pb.GetContactsRequest{UserID: req.UserId}, nil
}

func decodeGetContactsResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetContactsReply)

	return endpoint.GetContactsResponse{Contacts: res.Contacts, Err: str2err(res.Err)}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
}

// NewGRPCClient returns a UserService backed by a gRPC server at the other end
//  of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, options map[string][]grpcTrans.ClientOption) service.UsersService {
	return endpoint.Endpoints{
		CreateEndpoint:            makeCreateClient(conn, options[endpoint.Create]).Endpoint(),
		GetUserByEmailEndpoint:    makeGetUserByEmailClient(conn, options[endpoint.GetUserByEmail]).Endpoint(),
		GetUserByIDEndpoint:       makeGetUserByIDClient(conn, options[endpoint.GetUserById]).Endpoint(),
		GetUserByUUIDEndpoint:     makeGetUserByUUIDClient(conn, options[endpoint.GetUserByUUID]).Endpoint(),
		UpdateEmailEndpoint:       makeUpdateEmailClient(conn, options[endpoint.UpdateEmail]).Endpoint(),
		UpdateContactEndpoint:     makeUpdateContactClient(conn, options[endpoint.UpdateContact]).Endpoint(),
		UpdateUserProfileEndpoint: makeUpdateUserProfileClient(conn, options[endpoint.UpdateUserProfile]).Endpoint(),
		SetUserProfileEndpoint:    makeSetUserProfileClient(conn, options[endpoint.SetUserProfile]).Endpoint(),
		DeleteContactEndpoint:     makeDeleteContactClient(conn, options[endpoint.DeleteContact]).Endpoint(),
		CreateContactEndpoint:     makeCreateContactClient(conn, options[endpoint.CreateContact]).Endpoint(),
		GetContactsEndpoint:       makeGetContactsClient(conn, options[endpoint.GetContacts]).Endpoint(),
	}
}
