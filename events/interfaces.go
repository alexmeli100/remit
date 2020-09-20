package events

import (
	"context"
	eventpb "github.com/alexmeli100/remit/events/pb"
	paymentpb "github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	transferpb "github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	userpb "github.com/alexmeli100/remit/users/pkg/grpc/pb"
)

type UserEventSender interface {
	OnUserCreated(ctx context.Context, u *userpb.User) error
	OnPasswordReset(ctx context.Context, u *userpb.User) error
}

type TransferEventSender interface {
	OnTransferSucceded(ctx context.Context, t *transferpb.TransferResponse) error
}

type TransactionEventsSender interface {
	OnTransactionSucceded(ctx context.Context, t *paymentpb.Transaction) error
}

type PaymentEventSender interface {
	OnPaymentSucceded(ctx context.Context, paymentIntent string) error
}

type EventManager interface {
	UserEventSender
	TransferEventSender
	PaymentEventSender
	TransferEventSender
}

type UserEventHandler interface {
	OnUserCreated(ctx context.Context, data *eventpb.EventData) error
	OnPasswordReset(ctx context.Context, data *eventpb.EventData) error
}
