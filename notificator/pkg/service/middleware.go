package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(NotificatorService) NotificatorService

type loggingMiddleware struct {
	logger log.Logger
	next   NotificatorService
}

func (l loggingMiddleware) SendConfirmEmail(ctx context.Context, name, addr, link string) error {
	err := l.next.SendConfirmEmail(ctx, name, addr, link)

	defer func() {
		_ = l.logger.Log("method", "SendConfirmEmail", "name", name, "err", err)
	}()

	return err
}

func (l loggingMiddleware) SendPasswordResetEmail(ctx context.Context, addr, link string) error {
	err := l.next.SendPasswordResetEmail(ctx, addr, link)

	defer func() {
		_ = l.logger.Log("method", "SendPasswordReset", "err", err)
	}()

	return err
}

func (l loggingMiddleware) SendWelcomeEmail(ctx context.Context, name, addr string) error {
	err := l.next.SendWelcomeEmail(ctx, name, addr)

	defer func() {
		_ = l.logger.Log("method", "SendWelcomeEmail", "name", name, "err", err)
	}()

	return err
}

// LoggingMiddleware takes a logger as a dependency
// and returns a UsersService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next NotificatorService) NotificatorService {
		return &loggingMiddleware{logger, next}
	}
}
