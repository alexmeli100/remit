package grpc

import (
	"context"
	endpoint "github.com/alexmeli100/remit/transfer/pkg/endpoint"
	pb "github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"
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

	return &pb.TransferReply{Err: err2str(res.Err)}, nil
}

func (g *grpcServer) Transfer(ctx context1.Context, req *pb.TransferRequest) (*pb.TransferReply, error) {
	_, rep, err := g.transfer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.TransferReply), nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
