package grpc

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"

	"github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"
)

// makeCreateHandler creates the handler logic
func makeCreateHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.CreateEndpoint, decodeCreateRequest, encodeCreateResponse, options...)
}

// decodeCreateResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Create request.
func decodeCreateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(pb.CreateRequest)

	return endpoint.CreateRequest{User: req.User}, nil
}

// encodeCreateResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeCreateResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.CreateResponse)

	return pb.CreateReply{Err: res.Err.Error()}, nil
}
func (g *grpcServer) Create(ctx context1.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	_, rep, err := g.create.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.CreateReply), nil
}

// makeGetUserByIDHandler creates the handler logic
func makeGetUserByIDHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.GetUserByIDEndpoint, decodeGetUserByIDRequest, encodeGetUserByIDResponse, options...)
}

// decodeGetUserByIDResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetUserByID request.
func decodeGetUserByIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(pb.GetUserByIDRequest)

	return endpoint.GetUserByIDRequest{Id: req.Id}, nil
}

// encodeGetUserByIDResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGetUserByIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetUserByIDResponse)

	return pb.GetUserByIDReply{Err: res.Err.Error(), User: res.User}, nil
}
func (g *grpcServer) GetUserByID(ctx context1.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDReply, error) {
	_, rep, err := g.getUserByID.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.GetUserByIDReply), nil
}

// makeGetUserByEmailHandler creates the handler logic
func makeGetUserByEmailHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.GetUserByEmailEndpoint, decodeGetUserByEmailRequest, encodeGetUserByEmailResponse, options...)
}

// decodeGetUserByEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetUserByEmail request.
func decodeGetUserByEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(pb.GetUserByEmailRequest)

	return endpoint.GetUserByEmailRequest{Email: req.Email}, nil
}

// encodeGetUserByEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGetUserByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetUserByEmailResponse)

	return pb.GetUserByIDReply{User: res.User, Err: res.Err.Error()}, nil
}
func (g *grpcServer) GetUserByEmail(ctx context1.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailReply, error) {
	_, rep, err := g.getUserByEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetUserByEmailReply), nil
}

// makeUpdateEmailHandler creates the handler logic
func makeUpdateEmailHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.UpdateEmailEndpoint, decodeUpdateEmailRequest, encodeUpdateEmailResponse, options...)
}

// decodeUpdateEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain UpdateEmail request.
func decodeUpdateEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(pb.UpdateEmailRequest)

	return endpoint.UpdateEmailRequest{User: req.User}, nil
}

// encodeUpdateEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeUpdateEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.UpdateEmailResponse)

	return pb.UpdateEmailReply{Err: res.Err.Error()}, nil
}
func (g *grpcServer) UpdateEmail(ctx context1.Context, req *pb.UpdateEmailRequest) (*pb.UpdateEmailReply, error) {
	_, rep, err := g.updateEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateEmailReply), nil
}

// makeUpdatePasswordHandler creates the handler logic
func makeUpdatePasswordHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.UpdatePasswordEndpoint, decodeUpdatePasswordRequest, encodeUpdatePasswordResponse, options...)
}

// decodeUpdatePasswordResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain UpdatePassword request.

func decodeUpdatePasswordRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(pb.UpdatePasswordRequest)

	return endpoint.UpdatePasswordRequest{User: req.User}, nil
}

// encodeUpdatePasswordResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeUpdatePasswordResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.UpdatePasswordResponse)

	return pb.UpdatePasswordReply{Err: res.Err.Error()}, nil
}
func (g *grpcServer) UpdatePassword(ctx context1.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordReply, error) {
	_, rep, err := g.updatePassword.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdatePasswordReply), nil
}

// makeUpdateStatusHandler creates the handler logic
func makeUpdateStatusHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.UpdateStatusEndpoint, decodeUpdateStatusRequest, encodeUpdateStatusResponse, options...)
}

// decodeUpdateStatusResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain UpdateStatus request.
func decodeUpdateStatusRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(pb.UpdateStatusRequest)

	return endpoint.UpdateStatusRequest{User: req.User}, nil
}

// encodeUpdateStatusResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeUpdateStatusResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.UpdateStatusResponse)

	return pb.UpdateStatusReply{Err: res.Err.Error()}, nil
}
func (g *grpcServer) UpdateStatus(ctx context1.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusReply, error) {
	_, rep, err := g.updateStatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateStatusReply), nil
}
