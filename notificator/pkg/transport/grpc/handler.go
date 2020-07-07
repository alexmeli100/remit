package grpc

import (
	"context"
	"errors"
	"github.com/alexmeli100/remit/notificator/pkg/endpoint"
	"github.com/alexmeli100/remit/notificator/pkg/transport/grpc/pb"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// makeSendConfirmEmailHandler creates the handler logic
func makeSendConfirmEmailHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.SendConfirmEmailEndpoint, decodeSendConfirmEmailRequest, encodeSendConfirmEmailResponse, options...)
}

func makeSendConfirmEmailClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Notificator",
		endpoint.SendConfirmEmail,
		encodeSendConfirmEmailRequest,
		decodeSendConfirmEmailResponse,
		pb.SendConfirmEmailReply{},
		options...)
}

// decodeSendConfirmEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain SendConfirmEmail request.
func decodeSendConfirmEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SendConfirmEmailRequest)

	return &endpoint.SendConfirmEmailRequest{Link: req.Link, Name: req.Name, Addr: req.Addr}, nil
}

func decodeSendConfirmEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.SendConfirmEmailReply)

	return &endpoint.SendConfirmEmailResponse{Err: str2err(res.Err)}, nil
}

func encodeSendConfirmEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.SendConfirmEmailRequest)

	return &pb.SendConfirmEmailRequest{Name: req.Name, Link: req.Link, Addr: req.Addr}, nil
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

func makeSendPasswordResetEmailClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Notificator",
		endpoint.SendPasswordResetEmail,
		encodeSendPasswordResetEmailRequest,
		decodeSendPasswordResetEmailResponse,
		pb.SendPasswordResetEmailReply{},
		options...)
}

// decodeSendPasswordResetEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain SendPasswordResetEmail request.
func decodeSendPasswordResetEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SendPasswordResetEmailRequest)

	return &endpoint.SendPasswordResetEmailRequest{Addr: req.Addr, Link: req.Link}, nil
}

func decodeSendPasswordResetEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.SendPasswordResetEmailReply)

	return &endpoint.SendPasswordResetEmailResponse{Err: str2err(res.Err)}, nil
}

func encodeSendPasswordResetEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.SendPasswordResetEmailRequest)

	return &pb.SendPasswordResetEmailRequest{Link: req.Link, Addr: req.Addr}, nil
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

func makeSendWelcomeEmailClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Notificator",
		endpoint.SendWelcomeEmail,
		encodeSendWelcomeEmailRequest,
		decodeSendWelcomeEmailResponse,
		pb.SendWelcomeEmailReply{},
		options...)
}

func decodeSendWelcomeEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.SendWelcomeEmailReply)

	return &endpoint.SendWelcomeEmailResponse{Err: str2err(res.Err)}, nil
}

func encodeSendWelcomeEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*endpoint.SendWelcomeEmailRequest)

	return &pb.SendWelcomeEmailRequest{Name: req.Name, Addr: req.Addr}, nil
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

func str2err(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
