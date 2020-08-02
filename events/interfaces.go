package events

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
)

type UserEventManager interface {
	OnUserCreated(ctx context.Context, u *pb.User) error
	OnPasswordReset(ctx context.Context, u *pb.User) error
}

type EventManager interface {
	UserEventManager
}
