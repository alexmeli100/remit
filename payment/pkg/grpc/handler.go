package grpc

import (
	"context"
	"github.com/alexmeli100/remit/payment/pkg/endpoint"
	"github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	"github.com/go-kit/kit/transport/grpc"
)

// makeSaveCardHandler creates the handler logic
func makeSaveCardHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.SaveCardEndpoint, decodeSaveCardRequest, encodeSaveCardResponse, options...)
}

// decodeSaveCardRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain SaveCard request.
func decodeSaveCardRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SaveCardRequest)

	return endpoint.SaveCardRequest{Uid: req.Uid}, nil
}

// encodeSaveCardResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeSaveCardResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.SaveCardResponse)

	return &pb.SaveCardReply{Secret: res.Secret, Err: err2str(res.Err)}, nil
}
func (g *grpcServer) SaveCard(ctx context.Context, req *pb.SaveCardRequest) (*pb.SaveCardReply, error) {
	_, rep, err := g.saveCard.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.SaveCardReply), nil
}

// makeGetPaymentIntentSecretHandler creates the handler logic
func makeGetPaymentIntentSecretHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.GetPaymentIntentSecretEndpoint, decodeGetPaymentIntentSecretRequest, encodeGetPaymentIntentSecretResponse, options...)
}

// decodeGetPaymentIntentSecretRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetPaymentIntentSecret request.
func decodeGetPaymentIntentSecretRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetPaymentIntentSecretRequest)

	return endpoint.GetPaymentIntentSecretRequest{Req: req.Req}, nil
}

// encodeGetPaymentIntentSecretResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGetPaymentIntentSecretResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetPaymentIntentSecretResponse)

	return &pb.GetPaymentIntentSecretReply{Secret: res.Secret, Err: err2str(res.Err)}, nil
}
func (g *grpcServer) GetPaymentIntentSecret(ctx context.Context, req *pb.GetPaymentIntentSecretRequest) (*pb.GetPaymentIntentSecretReply, error) {
	_, rep, err := g.getPaymentIntentSecret.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.GetPaymentIntentSecretReply), nil
}

// makeGetPaymentIntentSecretHandler creates the handler logic
func makeCapturePaymentHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.CapturePaymentEndpoint, decodeCapturePaymentRequest, encodeCapturePaymentResponse, options...)
}

// decodeCapturePaymentRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetPaymentIntentSecret request.
func decodeCapturePaymentRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CapturePaymentRequest)

	return endpoint.CapturePaymentRequest{Amount: req.Amount, PI: req.Pi}, nil
}

// encodeCapturePaymentResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeCapturePaymentResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.CapturePaymentResponse)

	return &pb.CapturePaymentReply{Secret: res.Secret, Err: err2str(res.Err)}, nil
}

func (g *grpcServer) CapturePayment(ctx context.Context, req *pb.CapturePaymentRequest) (*pb.CapturePaymentReply, error) {
	_, res, err := g.capturePayment.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*pb.CapturePaymentReply), nil
}

func makeGetCustomerIDHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.GetCustomerIDEndpoint, decodeGetCustomerIDRequest, encodeGetCustomerIDResponse, options...)
}

func decodeGetCustomerIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetCustomerIDRequest)

	return endpoint.GetCustomerIDRequest{Uid: req.Uid}, nil
}

func encodeGetCustomerIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetCustomerIDResponse)

	return &pb.GetCustomerIDReply{CustomerID: res.CustomerID, Err: err2str(res.Err)}, nil
}

func (g *grpcServer) GetCustomerID(ctx context.Context, req *pb.GetCustomerIDRequest) (*pb.GetCustomerIDReply, error) {
	_, res, err := g.getCustomerID.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return res.(*pb.GetCustomerIDReply), nil
}

func makeCreateTransactionHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.CreateTransactionEndpoint, decodeCreateTransactionRequest, encodeCreateTransactionResponse, options...)
}

func decodeCreateTransactionRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CreateTransactionRequest)

	return endpoint.CreateTransactionRequest{Transaction: req.Transaction}, nil
}

func encodeCreateTransactionResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.CreateTransactionResponse)

	return &pb.CreateTransactionReply{Transaction: res.Transaction, Err: err2str(res.Err)}, nil
}

func (g *grpcServer) CreateTransaction(ctx context.Context, request *pb.CreateTransactionRequest) (*pb.CreateTransactionReply, error) {
	_, res, err := g.createTransaction.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return res.(*pb.CreateTransactionReply), nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
