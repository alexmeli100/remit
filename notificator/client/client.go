package client

import (
	"context"
	"errors"
	"github.com/alexmeli100/remit/notificator/pkg/endpoint"
	"github.com/alexmeli100/remit/notificator/pkg/service"
	"github.com/alexmeli100/remit/notificator/pkg/transport/grpc/pb"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

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

func decodeSendConfirmEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.SendConfirmEmailReply)

	return endpoint.SendConfirmEmailResponse{Err: str2err(res.Err)}, nil
}

func encodeSendConfirmEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.SendConfirmEmailRequest)

	return &pb.SendConfirmEmailRequest{Name: req.Name, Link: req.Link, Addr: req.Addr}, nil
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

func decodeSendPasswordResetEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.SendPasswordResetEmailReply)

	return endpoint.SendPasswordResetEmailResponse{Err: str2err(res.Err)}, nil
}

func encodeSendPasswordResetEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.SendPasswordResetEmailRequest)

	return &pb.SendPasswordResetEmailRequest{Link: req.Link, Addr: req.Addr}, nil
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

	return endpoint.SendWelcomeEmailResponse{Err: str2err(res.Err)}, nil
}

func encodeSendWelcomeEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.SendWelcomeEmailRequest)

	return &pb.SendWelcomeEmailRequest{Name: req.Name, Addr: req.Addr}, nil
}

func NewGRPCClient(conn *grpc.ClientConn, options map[string][]grpcTrans.ClientOption) service.NotificatorService {
	return endpoint.Endpoints{
		SendConfirmEmailEndpoint:       makeSendConfirmEmailClient(conn, options[endpoint.SendConfirmEmail]).Endpoint(),
		SendPasswordResetEmailEndpoint: makeSendPasswordResetEmailClient(conn, options[endpoint.SendPasswordResetEmail]).Endpoint(),
		SendWelcomeEmailEndpoint:       makeSendWelcomeEmailClient(conn, options[endpoint.SendWelcomeEmail]).Endpoint(),
	}
}

func str2err(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
}
