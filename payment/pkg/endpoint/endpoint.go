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
		s, err := s.SaveCard(ctx, req.Uid)

		return SaveCardResponse{
			Err:    err,
			Secret: s,
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
		s, err := s.GetPaymentIntentSecret(ctx, req.Req)
		return GetPaymentIntentSecretResponse{
			Err:    err,
			Secret: s,
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
