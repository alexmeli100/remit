package service

import (
	"context"
	"github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	"github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(TransferService) TransferService

type loggingMiddleware struct {
	logger log.Logger
	next   TransferService
}

func (l loggingMiddleware) Transfer(ctx context.Context, request *pb.TransferRequest) *pb.TransferResponse {
	res := l.next.Transfer(ctx, request)

	defer func() {
		l.logger.Log("method", "Transfer", "err", res.FailReason)
	}()

	return res
}

// LoggingMiddleware takes a logger as a dependency
// and returns a PaymentService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next TransferService) TransferService {
		return &loggingMiddleware{logger, next}
	}
}
