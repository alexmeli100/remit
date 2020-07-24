package grpc

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
)

// makeCreateHandler creates the handler logic
func makeCreateHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.CreateEndpoint, decodeCreateRequest, encodeCreateResponse, options...)
}

// decodeCreateRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Create request.
func decodeCreateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CreateRequest)

	return endpoint.CreateRequest{User: req.User}, nil
}

// encodeCreateResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeCreateResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.CreateResponse)

	return &pb.CreateReply{Err: err2str(res.Err)}, nil
}

func (g *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	_, rep, err := g.create.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.CreateReply), nil
}

// makeGetUserByIDHandler creates the handler logic
func makeGetUserByIDHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.GetUserByIDEndpoint, decodeGetUserByIDRequest, encodeGetUserByIDResponse, options...)
}

// decodeGetUserByIDResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetUserByID request.
func decodeGetUserByIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetUserByIDRequest)

	return endpoint.GetUserByIDRequest{Id: req.Id}, nil
}

// encodeGetUserByIDResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGetUserByIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetUserByIDResponse)

	return &pb.GetUserByIDReply{Err: err2str(res.Err), User: res.User}, nil
}

func (g *grpcServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDReply, error) {
	_, rep, err := g.getUserByID.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.GetUserByIDReply), nil
}

// makeGetUserByUUIDHandler creates the handler logic
func makeGetUserByUUIDHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.GetUserByUUIDEndpoint, decodeGetUserByUUIDRequest, encodeGetUserByUUIDResponse, options...)
}

// decodeGetUserByUUIDResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetUserByID request.
func decodeGetUserByUUIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetUserByUUIDRequest)

	return endpoint.GetUserByUUIDRequest{UUID: req.UUID}, nil
}

// encodeGetUserByUUIDResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGetUserByUUIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetUserByUUIDResponse)

	return &pb.GetUserByUUIDReply{Err: err2str(res.Err), User: res.User}, nil
}

func (g *grpcServer) GetUserByUUID(ctx context.Context, req *pb.GetUserByUUIDRequest) (*pb.GetUserByUUIDReply, error) {
	_, rep, err := g.getUserByUUID.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.GetUserByUUIDReply), nil
}

// makeGetUserByEmailHandler creates the handler logic
func makeGetUserByEmailHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.GetUserByEmailEndpoint, decodeGetUserByEmailRequest, encodeGetUserByEmailResponse, options...)
}

// decodeGetUserByEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetUserByEmail request.
func decodeGetUserByEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetUserByEmailRequest)

	return endpoint.GetUserByEmailRequest{Email: req.Email}, nil
}

// encodeGetUserByEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGetUserByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetUserByEmailResponse)

	return &pb.GetUserByEmailReply{User: res.User, Err: err2str(res.Err)}, nil
}

func (g *grpcServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailReply, error) {
	_, rep, err := g.getUserByEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetUserByEmailReply), nil
}

// makeUpdateEmailHandler creates the handler logic
func makeUpdateEmailHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.UpdateEmailEndpoint, decodeUpdateEmailRequest, encodeUpdateEmailResponse, options...)
}

// decodeUpdateEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain UpdateEmail request.
func decodeUpdateEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UpdateEmailRequest)

	return endpoint.UpdateEmailRequest{User: req.User}, nil
}

// encodeUpdateEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeUpdateEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.UpdateEmailResponse)

	return &pb.UpdateEmailReply{Err: err2str(res.Err)}, nil
}

func (g *grpcServer) UpdateEmail(ctx context.Context, req *pb.UpdateEmailRequest) (*pb.UpdateEmailReply, error) {
	_, rep, err := g.updateEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateEmailReply), nil
}

// makeUpdateStatusHandler creates the handler logic
func makeUpdateStatusHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.UpdateStatusEndpoint, decodeUpdateStatusRequest, encodeUpdateStatusResponse, options...)
}

// decodeUpdateStatusResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain UpdateStatus request.
func decodeUpdateStatusRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UpdateStatusRequest)

	return endpoint.UpdateStatusRequest{User: req.User}, nil
}

// encodeUpdateStatusResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeUpdateStatusResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.UpdateStatusResponse)

	return &pb.UpdateStatusReply{Err: err2str(res.Err)}, nil
}

func (g *grpcServer) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusReply, error) {
	_, rep, err := g.updateStatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateStatusReply), nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
