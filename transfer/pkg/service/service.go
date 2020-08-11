package service

import (
	"context"
	"github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
)

type SendMoney interface {
	SendTo(amount float64, recipient, currency string) error
}

// TransferService describes the service.
type TransferService interface {
	Transfer(ctx context.Context, request *pb.TransferRequest) error
}

//New returns a TransferService with all of the expected middleware wired in.
func New(svc TransferService, middleware []Middleware) TransferService {

	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
