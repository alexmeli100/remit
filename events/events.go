package events

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

type UserEventHandler func(ctx context.Context, user *pb.User) error

const (
	UserEvents = "user-events"
)

type UserEvent struct {
	EventKind string
	User      *pb.User
}

type EventSender struct {
	nats stan.Conn
}

func (e *EventSender) OnUserCreated(ctx context.Context, u *pb.User) error {
	b, err := encodeUserEvent(pb.UserCreated, u)

	if err != nil {
		return errors.Wrap(err, "proto marshal error")
	}

	return e.nats.Publish("user-events", b)
}

func (e *EventSender) OnPasswordReset(ctx context.Context, u *pb.User) error {
	b, err := encodeUserEvent(pb.UserPasswordReset, u)

	if err != nil {
		return errors.Wrap(err, "proto marshal error")
	}

	return e.nats.Publish("user-events", b)
}

func encodeUserEvent(kind pb.UserEventKind, u *pb.User) ([]byte, error) {
	event := &pb.UserEvent{
		Kind: kind,
		User: u,
	}

	return proto.Marshal(event)
}

func decodeUserEvent(data []byte) (*pb.UserEvent, error) {
	event := pb.UserEvent{}

	if err := proto.Unmarshal(data, &event); err != nil {
		return nil, errors.Wrap(err, "proto unmarshall error")
	}

	return &event, nil
}

func NewEventSender(conn stan.Conn) EventManager {
	return &EventSender{conn}
}

func Connect(url, clusterID, clientID string) (stan.Conn, error) {
	nc, err := nats.Connect(url)

	if err != nil {
		return nil, err
	}

	return stan.Connect(clusterID, clientID, stan.NatsConn(nc))
}

func ListenUserEvents(ctx context.Context, conn stan.Conn, queue string, handlers map[pb.UserEventKind]UserEventHandler, opts ...stan.SubscriptionOption) (chan error, error) {
	errc := make(chan error, 10)

	handler := func(msg *stan.Msg) {
		e, err := decodeUserEvent(msg.Data)

		if err != nil {
			errc <- err
			return
		}

		err = handlers[e.Kind](ctx, e.User)

		if err != nil {
			errc <- err
			return
		}

		if err = msg.Ack(); err != nil {
			errc <- err
		}
	}

	subs, err := conn.QueueSubscribe(UserEvents, queue, handler, opts...)

	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		subs.Unsubscribe()
		close(errc)
	}()

	return errc, nil
}

func ListenAllUserEvents(ctx context.Context, conn stan.Conn, queue string, sink UserEventManager, opts ...stan.SubscriptionOption) (chan error, error) {
	handlers := map[pb.UserEventKind]UserEventHandler{
		pb.UserCreated:       sink.OnUserCreated,
		pb.UserPasswordReset: sink.OnPasswordReset,
	}

	errc, err := ListenUserEvents(ctx, conn, queue, handlers, opts...)

	if err != nil {
		return nil, err
	}

	return errc, nil
}
