package grpc

import (
	"context"
	"github.com/alexmeli100/remit/transfer/pkg/endpoint"
	"github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	"github.com/alexmeli100/remit/transfer/pkg/service"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// NewGRPCClient returns an AddService backed by a gRPC server at the other end
//  of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, options map[string][]grpcTrans.ClientOption) service.TransferService {
	return endpoint.Endpoints{
		TransferEndpoint: makeTransferClient(conn, options[endpoint.Transfer]).Endpoint(),
	}
}

func makeTransferClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Transfer",
		"Transfer",
		encodeTransferRequest,
		decodeTransferResponse,
		pb.TransferResponse{},
		options...)
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
	res := r.(*pb.TransferResponse)

	return endpoint.TransferResponse{Res: res}, nil
}
