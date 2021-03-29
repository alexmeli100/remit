package service

import (
	"context"
	"go.uber.org/zap"
)

// Middleware describes a service middleware.
type Middleware func(TransferService) TransferService

type loggingMiddleware struct {
	logger *zap.Logger
	next   TransferService
}

func (l loggingMiddleware) Transfer(ctx context.Context, request *TransferRequest) *TransferResponse {
	res := l.next.Transfer(ctx, request)

	defer func() {
		l.logger.Info("Transfer", zap.String("err", res.FailReason))
	}()

	return res
}

// LoggingMiddleware takes a logger as a dependency
// and returns a TransferService Middleware.
func LoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next TransferService) TransferService {
		return &loggingMiddleware{logger, next}
	}
}
