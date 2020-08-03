package events

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/nats-io/nats.go"
)

type EventType int
type UserEventHandler func(ctx context.Context, user *pb.User) error

const (
	UserEvents = "user-events"
)

const (
	UserCreated       = "UserCreated"
	UserPasswordReset = "UserPasswordReset"
)

type UserEvent struct {
	EventKind string
	User      *pb.User
}

type EventSender struct {
	nats *nats.EncodedConn
}

func (e *EventSender) OnUserCreated(ctx context.Context, u *pb.User) error {
	return e.nats.Publish("user-events", UserEvent{UserCreated, u})
}

func (e *EventSender) OnPasswordReset(ctx context.Context, u *pb.User) error {
	return e.nats.Publish("user-events", UserEvent{UserPasswordReset, u})
}

func NewEventSender(conn *nats.EncodedConn) EventManager {
	return &EventSender{conn}
}

func Connect(url string) (*nats.EncodedConn, error) {
	conn, err := nats.Connect(url)

	if err != nil {
		return nil, err
	}

	return nats.NewEncodedConn(conn, nats.JSON_ENCODER)

}

func ListenUserEvents(ctx context.Context, conn *nats.EncodedConn, queue string, handlers map[string]UserEventHandler) (chan error, error) {
	errc := make(chan error, 10)

	subs, err := conn.QueueSubscribe(UserEvents, queue, func(e *UserEvent) {
		var err error
		err = handlers[e.EventKind](ctx, e.User)

		if err != nil {
			errc <- err
		}
	})

	if err != nil {
		return nil, err
	}

	go func() {
		defer subs.Unsubscribe()
		<-ctx.Done()
	}()

	return errc, nil
}

func ListenAllUserEvents(ctx context.Context, conn *nats.EncodedConn, queue string, sink UserEventManager) (chan error, error) {
	handlers := map[string]UserEventHandler{
		UserCreated:       sink.OnUserCreated,
		UserPasswordReset: sink.OnPasswordReset,
	}

	errc, err := ListenUserEvents(ctx, conn, queue, handlers)

	if err != nil {
		return nil, err
	}

	return errc, nil
}
