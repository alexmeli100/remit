package service

import (
	"context"
	"errors"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
)

var ErrUserNotFound = errors.New("user not found")

// UsersService describes the service.
type UsersService interface {
	Create(ctx context.Context, user *pb.User) error
	GetUserByID(ctx context.Context, id int64) (*pb.User, error)
	GetUserByUUID(ctx context.Context, uuid string) (*pb.User, error)
	GetUserByEmail(ctx context.Context, email string) (*pb.User, error)
	UpdateEmail(ctx context.Context, user *pb.User) error
	UpdateStatus(ctx context.Context, user *pb.User) error
}

// New returns a UsersService with all of the expected middleware wired in.
func New(svc UsersService, middleware []Middleware) UsersService {
	for _, m := range middleware {
		svc = m(svc)
	}

	return svc
}
