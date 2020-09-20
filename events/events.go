package events

import (
	"context"
	eventpb "github.com/alexmeli100/remit/events/pb"
	paymentpb "github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	transferpb "github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
	userpb "github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

type EventHandler func(ctx context.Context, data *eventpb.EventData) error
type Handlers map[eventpb.EventKind]EventHandler

const (
	UserEvents        = "user-events"
	PaymentEvents     = "payment-events"
	TransferEvents    = "transfer-events"
	TransactionEvents = "transaction-events"
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

func (e *EventSender) publishEvent(ctx context.Context, topic string, ev eventpb.EventKind, data *eventpb.EventData) error {
	event := &eventpb.Event{Event: ev, Payload: data}
	b, err := encodePbEvent(event)

	if err != nil {
		return errors.Wrap(err, "proto marshal error")
	}

	return e.nats.Publish(topic, b)
}

func (e *EventSender) OnUserCreated(ctx context.Context, u *userpb.User) error {
	data := &eventpb.EventData{
		Data: &eventpb.EventData_User{User: u},
	}

	return e.publishEvent(ctx, UserEvents, eventpb.UserCreated, data)
}

func (e *EventSender) OnPasswordReset(ctx context.Context, u *userpb.User) error {
	data := &eventpb.EventData{
		Data: &eventpb.EventData_User{User: u},
	}

	return e.publishEvent(ctx, UserEvents, eventpb.UserPasswordResert, data)
}

func (e *EventSender) OnTransferSucceded(ctx context.Context, t *transferpb.TransferResponse) error {
	data := &eventpb.EventData{
		Data: &eventpb.EventData_Transfer{Transfer: t},
	}

	return e.publishEvent(ctx, TransferEvents, eventpb.TransferSucceded, data)
}

func (e *EventSender) OnPaymentSucceded(ctx context.Context, paymentIntent string) error {
	data := &eventpb.EventData{
		Data: &eventpb.EventData_Intent{Intent: paymentIntent},
	}

	return e.publishEvent(ctx, PaymentEvents, eventpb.PaymentSucceded, data)
}

func (e *EventSender) onTransactionSucceded(ctx context.Context, t *paymentpb.Transaction) error {
	data := &eventpb.EventData{
		Data: &eventpb.EventData_Transaction{Transaction: t},
	}

	return e.publishEvent(ctx, TransactionEvents, eventpb.TransferSucceded, data)
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

func ListenEvents(ctx context.Context, topic, queue string, conn stan.Conn, handlers Handlers, opts ...stan.SubscriptionOption) (chan error, error) {
	errc := make(chan error, 1)

	handler := func(msg *stan.Msg) {
		e := &eventpb.Event{}
		err := decodePbEvent(msg.Data, e)

		if err != nil {
			errc <- err
			return
		}

		err = handlers[e.Event](ctx, e.Payload)

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

	subs, err := conn.QueueSubscribe(topic, queue, handler, opts...)

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

func ListenAllUserEvents(ctx context.Context, conn stan.Conn, queue string, sink UserEventHandler, opts ...stan.SubscriptionOption) (chan error, error) {
	handlers := map[eventpb.EventKind]EventHandler{
		eventpb.UserCreated:        sink.OnUserCreated,
		eventpb.UserPasswordResert: sink.OnPasswordReset,
	}

	errc, err := ListenEvents(ctx, UserEvents, queue, conn, handlers, opts...)

	if err != nil {
		return nil, err
	}

	return errc, nil
}
