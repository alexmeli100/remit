package service

import (
	"context"
	"github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	"github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(PaymentService) PaymentService

type loggingMiddleware struct {
	logger log.Logger
	next   PaymentService
}

func (l loggingMiddleware) GetCustomerID(ctx context.Context, uid string) (string, error) {
	c, err := l.next.GetCustomerID(ctx, uid)

	defer func() {
		l.logger.Log("method", "GetCustomerID", "err", err)
	}()

	return c, err
}

func (l loggingMiddleware) CreateTransaction(ctx context.Context, tr *pb.Transaction) (*pb.Transaction, error) {
	tr, err := l.next.CreateTransaction(ctx, tr)

	defer func() {
		l.logger.Log("method", "CreateTransaction", "err", err)
	}()

	return tr, err
}

// LoggingMiddleware takes a logger as a dependency
// and returns a PaymentService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next PaymentService) PaymentService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) SaveCard(ctx context.Context, uid string) (string, error) {
	s, err := l.next.SaveCard(ctx, uid)

	defer func() {
		l.logger.Log("method", "SaveCard", "uid", uid, "err", err)
	}()

	return s, err
}
func (l loggingMiddleware) GetPaymentIntentSecret(ctx context.Context, req *pb.PaymentRequest) (string, error) {
	s, err := l.next.GetPaymentIntentSecret(ctx, req)

	defer func() {
		l.logger.Log("method", "GetPaymentIntentSecret", "uid", req.Uid, "err", err)
	}()

	return s, err
}

func (l loggingMiddleware) CapturePayment(ctx context.Context, pi string, amount float64) (string, error) {
	s, err := l.next.CapturePayment(ctx, pi, amount)

	defer func() {
		l.logger.Log("method", "CapturePayment")
	}()

	return s, err
}
