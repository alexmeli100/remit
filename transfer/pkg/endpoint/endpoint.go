package endpoint

import (
	"context"
	"errors"
	"github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	"github.com/alexmeli100/remit/transfer/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
)

// TransferRequest collects the request parameters for the Transfer method.
type TransferRequest struct {
	Request *pb.TransferRequest `json:"request"`
}

// TransferResponse collects the response parameters for the Transfer method.
type TransferResponse struct {
	Res *pb.TransferResponse `json:"res"`
}

// MakeTransferEndpoint returns an endpoint that invokes Transfer on the service.
func MakeTransferEndpoint(s service.TransferService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tr := request.(TransferRequest)
		res := s.Transfer(ctx, tr.Request)

		return TransferResponse{Res: res}, nil
	}
}

// Failed implements Failer.
func (r TransferResponse) Failed() error {
	if r.Res.Status == "Failed" {
		return errors.New(r.Res.FailReason)
	}

	return nil
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Transfer implements Service. Primarily useful in a client.
func (e Endpoints) Transfer(ctx context.Context, req *pb.TransferRequest) *pb.TransferResponse {
	request := TransferRequest{Request: req}
	response, err := e.TransferEndpoint(ctx, request)

	if err != nil {
		return service.GetTransferResponse(req, err)
	}
	return response.(TransferResponse).Res
}
