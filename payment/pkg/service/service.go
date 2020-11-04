package service

import (
	"context"
	"github.com/alexmeli100/remit/payment/pkg/grpc/pb"
)

// PaymentService describes the service.
type PaymentService interface {
	SaveCard(ctx context.Context, uid string) (string, error)
	GetPaymentIntentSecret(ctx context.Context, req *pb.PaymentRequest) (string, error)
	CapturePayment(ctx context.Context, pi string, amount float64) (string, error)
	GetCustomerID(ctx context.Context, uid string) (string, error)
	CreateTransaction(ctx context.Context, tr *pb.Transaction) (*pb.Transaction, error)
}

// New returns a PaymentService with all of the expected middleware wired in.
func New(svc PaymentService, middleware []Middleware) PaymentService {
	for _, m := range middleware {
		svc = m(svc)
	}

	return svc
}
