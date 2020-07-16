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

// New returns a UserService backed by a gRPC server at the other end
//  of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func makeCreateClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "Create", encodeCreateRequest, decodeCreateResponse, pb.CreateReply{}, options...)
}

func encodeCreateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.CreateRequest)

	return &pb.CreateRequest{User: req.User}, nil
}

func decodeCreateResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.CreateReply)

	return &endpoint.CreateResponse{Err: str2err(res.Err)}, nil
}

func makeGetUserByIDClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "GetUserByID", encodeGetUserByIDRequest, decodeGetUserByIDResponse, pb.GetUserByIDReply{}, options...)
}

func decodeGetUserByIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetUserByIDReply)

	return &endpoint.GetUserByIDResponse{User: res.User, Err: str2err(res.Err)}, nil
}

func encodeGetUserByIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.GetUserByIDRequest)

	return &pb.GetUserByIDRequest{Id: req.Id}, nil
}

func makeGetUserByUUIDClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "GetUserByUUID", encodeGetUserByUUIDRequest, decodeGetUserByUUIDResponse, pb.GetUserByUUIDReply{}, options...)
}

func decodeGetUserByUUIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetUserByIDReply)

	return &endpoint.GetUserByIDResponse{User: res.User, Err: str2err(res.Err)}, nil
}

func encodeGetUserByUUIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.GetUserByUUIDRequest)

	return &pb.GetUserByUUIDRequest{UUID: req.UUID}, nil
}

func makeGetUserByEmailClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "GetUserByEmail", encodeGetUserByEmailRequest, decodeGetUserByEmailResponse, pb.GetUserByEmailReply{}, options...)
}

func decodeGetUserByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetUserByEmailReply)

	return &endpoint.GetUserByEmailResponse{User: res.User, Err: str2err(res.Err)}, nil
}

func encodeGetUserByEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.GetUserByEmailRequest)

	return &pb.GetUserByEmailRequest{Email: req.Email}, nil
}

func makeUpdateEmailClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "UpdateEmail", encodeUpdateEmailRequest, decodeUpdateEmailResponse, pb.UpdateEmailReply{}, options...)
}

func decodeUpdateEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.UpdateEmailReply)

	return &endpoint.UpdateEmailResponse{Err: str2err(res.Err)}, nil
}

func encodeUpdateEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.UpdateEmailRequest)

	return &pb.UpdateEmailRequest{User: req.User}, nil
}

func makeUpdateStatusClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(conn, "pb.Users", "UpdateStatus", encodeUpdateStatusRequest, decodeUpdateStatusResponse, pb.UpdateStatusReply{}, options...)
}

func decodeUpdateStatusResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.UpdateStatusReply)

	return &endpoint.UpdateStatusResponse{Err: str2err(res.Err)}, nil
}

func encodeUpdateStatusRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.UpdateStatusRequest)

	return &pb.UpdateStatusRequest{User: req.User}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
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
