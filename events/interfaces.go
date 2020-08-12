package events

import (
	"context"
	transferpb "github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
)

type UserEventManager interface {
	OnUserCreated(ctx context.Context, u *pb.User) error
	OnPasswordReset(ctx context.Context, u *pb.User) error
}

type TransferEventManager interface {
	OnTransferSucceded(ctx context.Context, t *transferpb.TransferRequest) error
}

type PaymentEventManager interface {
	OnPaymentSucceded(ctx context.Context, paymentIntent string) error
}

type EventManager interface {
	UserEventManager
	TransferEventManager
	PaymentEventManager
}
