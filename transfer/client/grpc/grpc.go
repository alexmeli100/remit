package grpc

import (
	"context"
	"errors"
	"github.com/alexmeli100/remit/transfer/pkg/endpoint"
	pb "github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	service "github.com/alexmeli100/remit/transfer/pkg/service"
	kitEndpoint "github.com/go-kit/kit/endpoint"
	grpc1 "github.com/go-kit/kit/transport/grpc"
	grpc "google.golang.org/grpc"
)

// New returns an AddService backed by a gRPC server at the other end
//  of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func New(conn *grpc.ClientConn, options map[string][]grpc1.ClientOption) (service.TransferService, error) {
	var transferEndpoint kitEndpoint.Endpoint
	{
		transferEndpoint = grpc1.NewClient(conn, "pb.Transfer", "Transfer", encodeTransferRequest, decodeTransferResponse, pb.TransferReply{}, options["Transfer"]...).Endpoint()
	}

	return endpoint.Endpoints{TransferEndpoint: transferEndpoint}, nil
}

// encodeTransferRequest is a transport/grpc.EncodeRequestFunc that converts a
//  user-domain Transfer request to a gRPC request.
func encodeTransferRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.TransferRequest)

	return req.Request, nil
}

// decodeTransferResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response.
func decodeTransferResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.TransferReply)

	return endpoint.TransferResponse{Err: str2err(res.Err)}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
}
