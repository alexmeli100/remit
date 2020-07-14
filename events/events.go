package events

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/nats-io/nats.go"
)

type EventType int

const (
	UserCreated EventType = iota
	UserPasswordReset
)

type UserEvent struct {
	EventKind EventType
	User      *pb.User
}

type EventSender struct {
	hostname string
	nats     *nats.EncodedConn
}

func (e *EventSender) OnUserCreated(ctx context.Context, u *pb.User) error {
	return e.nats.Publish("user-events", UserEvent{UserCreated, u})
}

func (e *EventSender) OnPasswordReset(ctx context.Context, u *pb.User) error {
	return e.nats.Publish("user-events", UserEvent{UserPasswordReset, u})
}

func NewEventSender(instance string) (EventManager, error) {
	conn, err := connect(instance)

	if err != nil {
		return nil, err
	}

	return &EventSender{instance, conn}, nil
}

func connect(url string) (*nats.EncodedConn, error) {
	conn, err := nats.Connect(url)

	if err != nil {
		return nil, err
	}

	return nats.NewEncodedConn(conn, nats.JSON_ENCODER)

}