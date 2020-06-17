package service

import (
	"context"
	pb "github.com/alexmeli100/remit/users/pkg/grpc/pb"

	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(UsersService) UsersService

type loggingMiddleware struct {
	logger log.Logger
	next   UsersService
}

func (l loggingMiddleware) Create(ctx context.Context, user *pb.User) error {
	err := l.next.Create(ctx, user)

	defer func() {
		_ = l.logger.Log("method", "Create", "user", user, "err", err)
	}()

	return err
}

func (l loggingMiddleware) GetUserByID(ctx context.Context, id int64) (*pb.User, error) {
	user, err := l.next.GetUserByID(ctx, id)

	defer func() {
		_ = l.logger.Log("method", "GetUserById", "user", user, "err", err)
	}()

	return user, err
}

func (l loggingMiddleware) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	user, err := l.next.GetUserByEmail(ctx, email)

	defer func() {
		_ = l.logger.Log("method", "GetUserByEmail", "user", user, "err", err)
	}()

	return user, err
}

func (l loggingMiddleware) UpdateEmail(ctx context.Context, user *pb.User) error {
	err := l.next.UpdateEmail(ctx, user)

	defer func() {
		_ = l.logger.Log("method", "UpdateEmail", "user", user, "err", err)
	}()

	return err
}

func (l loggingMiddleware) UpdatePassword(ctx context.Context, user *pb.User) error {
	err := l.next.UpdatePassword(ctx, user)

	defer func() {
		_ = l.logger.Log("method", "UpdatePassword", "user", user, "err", err)
	}()

	return err
}

func (l loggingMiddleware) UpdateStatus(ctx context.Context, user *pb.User) error {
	err := l.next.UpdateStatus(ctx, user)

	defer func() {
		_ = l.logger.Log("method", "Updatestatus", "user", user, "err", err)
	}()

	return err
}

// LoggingMiddleware takes a logger as a dependency
// and returns a UsersService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next UsersService) UsersService {
		return &loggingMiddleware{logger, next}
	}
}
