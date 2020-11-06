package endpoint

import (
	"context"
	"github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	"github.com/alexmeli100/remit/payment/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

// SaveCardRequest collects the request parameters for the SaveCard method.
type SaveCardRequest struct {
	Uid string `json:"uid"`
}

// SaveCardResponse collects the response parameters for the SaveCard method.
type SaveCardResponse struct {
	Secret string `json:"secret"`
	Err    error  `json:"err"`
}

// MakeSaveCardEndpoint returns an endpoint that invokes SaveCard on the service.
func MakeSaveCardEndpoint(s service.PaymentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SaveCardRequest)
		secret, err := s.SaveCard(ctx, req.Uid)

		return SaveCardResponse{
			Err:    err,
			Secret: secret,
		}, nil
	}
}

// Failed implements Failer.
func (r SaveCardResponse) Failed() error {
	return r.Err
}

// CapturePaymentRequest collects the request parameters for the CapturePayment method
type CapturePaymentRequest struct {
	PI     string  `json:"pi"`
	Amount float64 `json:"amount"`
}

// CapturePaymentResponse collects the response paramters for the CapturePayment method
type CapturePaymentResponse struct {
	Secret string `json:"secret"`
	Err    error  `json:"err"`
}

func MakeCapturePaymentEndpoint(s service.PaymentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CapturePaymentRequest)
		s, err := s.CapturePayment(ctx, req.PI, req.Amount)

		return CapturePaymentResponse{Secret: s, Err: err}, nil
	}
}

func (r CapturePaymentResponse) Failed() error {
	return r.Err
}

type GetCustomerIDRequest struct {
	Uid string `json:"uid"`
}

type GetCustomerIDResponse struct {
	CustomerID string `json:"customerID"`
	Err        error  `json:"err"`
}

func MakeGetCustomerIDEndpoint(s service.PaymentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetCustomerIDRequest)
		c, err := s.GetCustomerID(ctx, req.Uid)

		return GetCustomerIDResponse{CustomerID: c, Err: err}, nil
	}
}

func (r GetCustomerIDResponse) Failed() error {
	return r.Err
}

type CreateTransactionRequest struct {
	Transaction *pb.Transaction `json:"transaction"`
}

type CreateTransactionResponse struct {
	Transaction *pb.Transaction `json:"transaction"`
	Err         error           `json:"err"`
}

func MakeCreateTransactionEndpoint(s service.PaymentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateTransactionRequest)
		tr, err := s.CreateTransaction(ctx, req.Transaction)

		return CreateTransactionResponse{Transaction: tr, Err: err}, nil
	}
}

func (r CreateTransactionResponse) Failed() error {
	return r.Err
}

// GetPaymentIntentSecretRequest collects the request parameters for the GetPaymentIntentSecret method.
type GetPaymentIntentSecretRequest struct {
	Req *pb.PaymentRequest `json:"req"`
}

// GetPaymentIntentSecretResponse collects the response parameters for the GetPaymentIntentSecret method.
type GetPaymentIntentSecretResponse struct {
	Secret string `json:"secret"`
	Err    error  `json:"err"`
}

// MakeGetPaymentIntentSecretEndpoint returns an endpoint that invokes GetPaymentIntentSecret on the service.
func MakeGetPaymentIntentSecretEndpoint(s service.PaymentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetPaymentIntentSecretRequest)
		secret, err := s.GetPaymentIntentSecret(ctx, req.Req)

		return GetPaymentIntentSecretResponse{
			Err:    err,
			Secret: secret,
		}, nil
	}
}

// Failed implements Failer.
func (r GetPaymentIntentSecretResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

type GetTransactionsRequest struct {
	Uid string `json:"uid"`
}

type GetTransactionsResponse struct {
	Transactions []*pb.Transaction
	Err          error
}

func MakeGetTransactionsEndpoint(s service.PaymentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetCustomerIDRequest)
		trs, err := s.GetTransactions(ctx, req.Uid)

		if err != nil {
			return nil, err
		}

		res := GetTransactionsResponse{Transactions: trs, Err: err}
		return res, nil
	}
}

// SaveCard implements Service. Primarily useful in a client.
func (e Endpoints) SaveCard(ctx context.Context, uid string) (string, error) {
	request := SaveCardRequest{Uid: uid}
	response, err := e.SaveCardEndpoint(ctx, request)
	if err != nil {
		return "", err
	}
	return response.(SaveCardResponse).Secret, response.(SaveCardResponse).Err
}

// GetPaymentIntentSecret implements Service. Primarily useful in a client.
func (e Endpoints) GetPaymentIntentSecret(ctx context.Context, req *pb.PaymentRequest) (string, error) {
	request := GetPaymentIntentSecretRequest{Req: req}
	response, err := e.GetPaymentIntentSecretEndpoint(ctx, request)

	if err != nil {
		return "", err
	}
	return response.(GetPaymentIntentSecretResponse).Secret, response.(GetPaymentIntentSecretResponse).Err
}

func (e Endpoints) CapturePayment(ctx context.Context, pi string, amount float64) (string, error) {
	request := CapturePaymentRequest{PI: pi, Amount: amount}
	response, err := e.CapturePaymentEndpoint(ctx, request)

	if err != nil {
		return "", err
	}

	res := response.(CapturePaymentResponse)
	return res.Secret, res.Err
}

func (e Endpoints) GetCustomerID(ctx context.Context, uid string) (string, error) {
	request := GetCustomerIDRequest{Uid: uid}
	response, err := e.GetCustomerIDEndpoint(ctx, request)

	if err != nil {
		return "", err
	}

	res := response.(GetCustomerIDResponse)
	return res.CustomerID, res.Err
}

func (e Endpoints) CreateTransaction(ctx context.Context, tr *pb.Transaction) (*pb.Transaction, error) {
	request := CreateTransactionRequest{Transaction: tr}
	response, err := e.CreateTransactionEndpoint(ctx, request)

	if err != nil {
		return nil, err
	}

	res := response.(CreateTransactionResponse)
	return res.Transaction, res.Err
}

func (e Endpoints) GetTransactions(ctx context.Context, uid string) ([]*pb.Transaction, error) {
	request := GetTransactionsRequest{Uid: uid}
	response, err := e.GetTransactionsEndpoint(ctx, request)

	if err != nil {
		return nil, err
	}

	res := response.(GetTransactionsResponse)

	return res.Transactions, res.Err
}
