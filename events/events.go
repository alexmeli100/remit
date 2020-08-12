package events

import (
	"context"
	paymentpb "github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	transferpb "github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

type UserEventHandler func(ctx context.Context, user *pb.User) error
type TransferEventHandler func(ctx context.Context, t *transferpb.TransferRequest) error
type PaymentEventHandler func(ctx context.Context, intent string) error

const (
	UserEvents     = "user-events"
	PaymentEvents  = "payment-events"
	TransferEvents = "transfer-events"
)

type ErrorShouldAck struct {
	Err string
}

func (e *ErrorShouldAck) Error() string {
	return e.Err
}

type EventSender struct {
	nats stan.Conn
}

func (e *EventSender) OnUserCreated(ctx context.Context, u *pb.User) error {
	event := &pb.UserEvent{User: u, Kind: pb.UserCreated}
	b, err := encodePbEvent(event)

	if err != nil {
		return errors.Wrap(err, "proto marshal error")
	}

	return e.nats.Publish(UserEvents, b)
}

func (e *EventSender) OnPasswordReset(ctx context.Context, u *pb.User) error {
	event := &pb.UserEvent{User: u, Kind: pb.UserPasswordReset}
	b, err := encodePbEvent(event)

	if err != nil {
		return errors.Wrap(err, "proto marshal error")
	}

	return e.nats.Publish(UserEvents, b)
}

func (e *EventSender) OnTransferSucceded(ctx context.Context, t *transferpb.TransferRequest) error {
	event := &transferpb.TransferEvent{Request: t, Kind: transferpb.TransferSucceded}
	b, err := encodePbEvent(event)

	if err != nil {
		return errors.Wrap(err, "proto marshal error")
	}

	return e.nats.Publish(TransferEvents, b)
}

func (e *EventSender) OnPaymentSucceded(ctx context.Context, paymentIntent string) error {
	event := &paymentpb.PaymentEvent{Kind: paymentpb.PaymentSucceded, Intent: paymentIntent}
	b, err := encodePbEvent(event)

	if err != nil {
		return errors.Wrap(err, "proto marshal error")
	}

	return e.nats.Publish(TransferEvents, b)
}

func encodePbEvent(src proto.Message) ([]byte, error) {
	return proto.Marshal(src)
}

func decodePbEvent(data []byte, dst proto.Message) error {
	if err := proto.Unmarshal(data, dst); err != nil {
		return errors.Wrap(err, "proto unmarshall error")
	}

	return nil
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
		var e pb.UserEvent
		err := decodePbEvent(msg.Data, &e)

		if err != nil {
			errc <- err
			return
		}

		err = handlers[e.Kind](ctx, e.User)

		if err != nil {
			errc <- err

			var e *ErrorShouldAck
			if errors.As(err, e) {
				sendAck(msg, errc)
			}

			return
		}

		sendAck(msg, errc)
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

func sendAck(msg *stan.Msg, errc chan error) {
	if err := msg.Ack(); err != nil {
		errc <- err
	}
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
