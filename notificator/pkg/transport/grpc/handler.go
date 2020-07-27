package grpc

import (
	"context"
	"github.com/alexmeli100/remit/notificator/pkg/endpoint"
	"github.com/alexmeli100/remit/notificator/pkg/transport/grpc/pb"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
)

// makeSendConfirmEmailHandler creates the handler logic
func makeSendConfirmEmailHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.SendConfirmEmailEndpoint, decodeSendConfirmEmailRequest, encodeSendConfirmEmailResponse, options...)
}

// decodeSendConfirmEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain SendConfirmEmail request.
func decodeSendConfirmEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SendConfirmEmailRequest)

	return &endpoint.SendConfirmEmailRequest{Link: req.Link, Name: req.Name, Addr: req.Addr}, nil
}

// encodeSendConfirmEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeSendConfirmEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*endpoint.SendConfirmEmailResponse)

	return &pb.SendConfirmEmailReply{Err: err2str(res.Err)}, nil
}

func (g *grpcServer) SendConfirmEmail(ctx context.Context, req *pb.SendConfirmEmailRequest) (*pb.SendConfirmEmailReply, error) {
	_, rep, err := g.sendConfirmEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SendConfirmEmailReply), nil
}

// makeSendPasswordResetEmailHandler creates the handler logic
func makeSendPasswordResetEmailHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.SendPasswordResetEmailEndpoint, decodeSendPasswordResetEmailRequest, encodeSendPasswordResetEmailResponse, options...)
}

// decodeSendPasswordResetEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain SendPasswordResetEmail request.
func decodeSendPasswordResetEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SendPasswordResetEmailRequest)

	return &endpoint.SendPasswordResetEmailRequest{Addr: req.Addr, Link: req.Link}, nil
}

// encodeSendPasswordResetEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeSendPasswordResetEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*endpoint.SendPasswordResetEmailResponse)

	return &pb.SendPasswordResetEmailReply{Err: err2str(res.Err)}, nil
}
func (g *grpcServer) SendPasswordResetEmail(ctx context.Context, req *pb.SendPasswordResetEmailRequest) (*pb.SendPasswordResetEmailReply, error) {
	_, rep, err := g.sendPasswordResetEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SendPasswordResetEmailReply), nil
}

// makeSendWelcomeEmailHandler creates the handler logic
func makeSendWelcomeEmailHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.SendWelcomeEmailEndpoint, decodeSendWelcomeEmailRequest, encodeSendWelcomeEmailResponse, options...)
}

// decodeSendWelcomeEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain SendWelcomeEmail request.
func decodeSendWelcomeEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SendWelcomeEmailRequest)

	return &endpoint.SendWelcomeEmailRequest{Name: req.Name, Addr: req.Addr}, nil
}

// encodeSendWelcomeEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeSendWelcomeEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*endpoint.SendWelcomeEmailResponse)

	return &pb.SendWelcomeEmailReply{Err: err2str(res.Err)}, nil
}
func (g *grpcServer) SendWelcomeEmail(ctx context.Context, req *pb.SendWelcomeEmailRequest) (*pb.SendWelcomeEmailReply, error) {
	_, rep, err := g.sendWelcomeEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SendWelcomeEmailReply), nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
