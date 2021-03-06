package grpc

import (
	"context"
	"errors"
	endpoint "github.com/alexmeli100/remit/payment/pkg/endpoint"
	pb "github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	service "github.com/alexmeli100/remit/payment/pkg/service"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	grpc "google.golang.org/grpc"
)

// New returns an AddService backed by a gRPC server at the other end
//  of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func New(conn *grpc.ClientConn, options map[string][]grpcTrans.ClientOption) service.PaymentService {
	return endpoint.Endpoints{
		CapturePaymentEndpoint:         makeCapturePaymentClient(conn, options[endpoint.CapturePayment]).Endpoint(),
		GetPaymentIntentSecretEndpoint: makeGetPaymentSecretClient(conn, options[endpoint.CapturePayment]).Endpoint(),
		SaveCardEndpoint:               makeSaveCardClient(conn, options[endpoint.SaveCard]).Endpoint(),
		GetCustomerIDEndpoint:          makeGetCustomerIDClient(conn, options[endpoint.GetCustomerID]).Endpoint(),
		CreateTransactionEndpoint:      makeCreateTransactionClient(conn, options[endpoint.CreateTransaction]).Endpoint(),
		GetTransactionsEndpoint:        makeGetTransactionsClient(conn, options[endpoint.GetTransactions]).Endpoint(),
	}
}

func makeSaveCardClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Payment",
		endpoint.SaveCard,
		encodeSaveCardRequest,
		decodeSaveCardResponse,
		pb.SaveCardReply{},
		options...)

}

// encodeSaveCardRequest is a transport/grpc.EncodeRequestFunc that converts a
//  user-domain SaveCard request to a gRPC request.
func encodeSaveCardRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.SaveCardRequest)

	return &pb.SaveCardRequest{Uid: req.Uid}, nil
}

// decodeSaveCardResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response.
func decodeSaveCardResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.SaveCardReply)

	return endpoint.SaveCardResponse{Secret: res.Secret, Err: str2err(res.Err)}, nil
}

func makeGetPaymentSecretClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Payment",
		endpoint.GetPaymentIntentSecret,
		encodeGetPaymentIntentSecretRequest,
		decodeGetPaymentIntentSecretResponse,
		pb.GetPaymentIntentSecretReply{},
		options...)

}

// encodeGetPaymentIntentSecretRequest is a transport/grpc.EncodeRequestFunc that converts a
//  user-domain GetPaymentIntentSecret request to a gRPC request.
func encodeGetPaymentIntentSecretRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.GetPaymentIntentSecretRequest)

	return &pb.GetPaymentIntentSecretRequest{Req: req.Req}, nil
}

// decodeGetPaymentIntentSecretResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response.
func decodeGetPaymentIntentSecretResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetPaymentIntentSecretReply)

	return endpoint.GetPaymentIntentSecretResponse{Secret: res.Secret, Err: str2err(res.Err)}, nil
}

func makeCapturePaymentClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Payment",
		endpoint.CapturePayment,
		encodeCapturePaymentRequest,
		decodeCapturePaymentResponse,
		pb.CapturePaymentReply{},
		options...)

}

// encodeCapturePaymentRequest is a transport/grpc.EncodeRequestFunc that converts a
//  user-domain CapturePayment request to a gRPC request.
func encodeCapturePaymentRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.CapturePaymentRequest)

	return &pb.CapturePaymentRequest{Pi: req.PI, Amount: req.Amount}, nil
}

// decodeCapturePaymentResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response.
func decodeCapturePaymentResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.CapturePaymentReply)

	return endpoint.CapturePaymentResponse{Secret: res.Secret, Err: str2err(res.Err)}, nil
}

func makeGetCustomerIDClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Payment",
		endpoint.GetCustomerID,
		encodeGetCustomerIDRequest,
		decodeGetCustomerIDResponse,
		pb.GetCustomerIDReply{},
		options...)
}

func decodeGetCustomerIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetCustomerIDReply)

	return endpoint.GetCustomerIDResponse{Err: str2err(res.Err), CustomerID: res.CustomerID}, nil
}

func encodeGetCustomerIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.GetCustomerIDRequest)

	return &pb.GetCustomerIDRequest{Uid: req.Uid}, nil
}

func makeCreateTransactionClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Payment",
		endpoint.CreateTransaction,
		encodeCreateTransactionRequest,
		decodeCreateTransactionResponse,
		pb.CreateTransactionReply{},
		options...)
}

func decodeCreateTransactionResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.CreateTransactionReply)

	return endpoint.CreateTransactionResponse{Err: str2err(res.Err), Transaction: res.Transaction}, nil
}

func encodeCreateTransactionRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.CreateTransactionRequest)

	return &pb.CreateTransactionRequest{Transaction: req.Transaction}, nil
}

func makeGetTransactionsClient(conn *grpc.ClientConn, options []grpcTrans.ClientOption) *grpcTrans.Client {
	return grpcTrans.NewClient(
		conn,
		"pb.Payment",
		endpoint.GetTransactions,
		encodeGetTransactionsRequest,
		decodeGetTransactionsResponse,
		pb.GetTransactionsReply{},
		options...)
}

func decodeGetTransactionsResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.GetTransactionsReply)

	return endpoint.GetTransactionsResponse{Err: str2err(res.Err), Transactions: res.Transactions}, nil
}

func encodeGetTransactionsRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.GetTransactionsRequest)

	return &pb.GetTransactionsRequest{Uid: req.Uid}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
}
