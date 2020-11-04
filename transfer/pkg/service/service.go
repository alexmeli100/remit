package service

import (
	"context"
	"github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
)

type SendMoney interface {
	SendTo(amount int, recipient, currency string) error
}

// TransferService describes the service.
type TransferService interface {
	Transfer(ctx context.Context, request *pb.TransferRequest) *pb.TransferResponse
}

//New returns a TransferService with all of the expected middleware wired in.
func New(svc TransferService, middleware []Middleware) TransferService {

	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
