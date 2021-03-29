package service

import (
	"context"
	"go.uber.org/zap"
)

// Middleware describes a service middleware.
type Middleware func(NotificatorService) NotificatorService

type loggingMiddleware struct {
	logger *zap.Logger
	next   NotificatorService
}

func (l loggingMiddleware) SendConfirmEmail(ctx context.Context, name, addr, link string) error {
	err := l.next.SendConfirmEmail(ctx, name, addr, link)

	defer func() {
		l.logger.Info("SendConfirmEmail", zap.String("name", name), zap.Error(err))
	}()

	return err
}

func (l loggingMiddleware) SendPasswordResetEmail(ctx context.Context, addr, link string) error {
	err := l.next.SendPasswordResetEmail(ctx, addr, link)

	defer func() {
		l.logger.Info("SendPasswordResetEmail", zap.Error(err))
	}()

	return err
}

func (l loggingMiddleware) SendWelcomeEmail(ctx context.Context, name, addr string) error {
	err := l.next.SendWelcomeEmail(ctx, name, addr)

	defer func() {
		l.logger.Info("SendWelcomeEmail", zap.String("name", name), zap.Error(err))
	}()

	return err
}

// LoggingMiddleware takes a logger as a dependency
// and returns a UsersService Middleware.
func LoggingMiddleware(logger *zap.Logger) Middleware {
	return func(next NotificatorService) NotificatorService {
		return &loggingMiddleware{logger, next}
	}
}
