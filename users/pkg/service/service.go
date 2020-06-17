package service

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/alexmeli100/remit/users/pkg/service/postgres"
)

// UsersService describes the service.
type UsersService interface {
	Create(ctx context.Context, user *pb.User) error
	GetUserByID(ctx context.Context, id int64) (*pb.User, error)
	GetUserByEmail(ctx context.Context, email string) (*pb.User, error)
	UpdateEmail(ctx context.Context, user *pb.User) error
	UpdatePassword(ctx context.Context, user *pb.User) error
	UpdateStatus(ctx context.Context, user *pb.User) error
}

// New returns a UsersService with all of the expected middleware wired in.
func New(middleware []Middleware) UsersService {
	svc := postgres.NewPostgService()

	for _, m := range middleware {
		svc = m(svc)
	}

	return svc
}
