package grpc

import (
	"context"
	endpoint "github.com/alexmeli100/remit/transfer/pkg/endpoint"
	pb "github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	"github.com/alexmeli100/remit/transfer/pkg/service"
	grpc "github.com/go-kit/kit/transport/grpc"
)

// makeTransferHandler creates the handler logic
func makeTransferHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.TransferEndpoint, decodeTransferRequest, encodeTransferResponse, options...)
}

// decodeTransferRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Transfer request.
func decodeTransferRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.TransferRequest)

	return endpoint.TransferRequest{Request: req}, nil
}

// encodeTransferResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeTransferResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.TransferResponse)

	return res.Res, nil
}

func (g *grpcServer) Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	_, rep, err := g.transfer.ServeGRPC(ctx, req)
	if err != nil {
		return service.GetTransferResponse(req, err), err
	}
	return rep.(*pb.TransferResponse), nil
}
