package service

import (
	"context"
)

// NotificatorService describes the service.
type NotificatorService interface {
	SendConfirmEmail(ctx context.Context, name, addr, link string) error
	SendPasswordResetEmail(ctx context.Context, addr, link string) error
	SendWelcomeEmail(ctx context.Context, name, addr string) error
}

// New returns a NotificatorService with all of the expected middleware wired in.
func New(svc NotificatorService, middleware []Middleware) NotificatorService {

	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
