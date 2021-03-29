package service

import (
	"context"
	"go.uber.org/zap"
)

// Middleware describes a service middleware.
type Middleware func(PaymentService) PaymentService

type loggingMiddleware struct {
	logger *zap.Logger
	next   PaymentService
}

func (l loggingMiddleware) GetTransactions(ctx context.Context, uid string) ([]*Transaction, error) {
	trs, err := l.next.GetTransactions(ctx, uid)

	defer func() {
		l.logger.Info("GetTransactions", zap.Error(err))
	}()

	return trs, err
}

func (l loggingMiddleware) CreateTransaction(ctx context.Context, tr *Transaction) (*Transaction, error) {
	tr, err := l.next.CreateTransaction(ctx, tr)

	defer func() {
		l.logger.Info("CreateTransaction", zap.Error(err))
	}()

	return tr, err
}

// LoggingMiddleware takes a logger as a dependency
// and returns a PaymentService Middleware.
func LoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next PaymentService) PaymentService {
		return &loggingMiddleware{logger, next}
	}

}
