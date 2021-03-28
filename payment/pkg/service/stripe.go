package service

import (
	"context"
	"fmt"
	"github.com/alexmeli100/remit/events"
	eventpb "github.com/alexmeli100/remit/events/pb"
	"github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	userPb "github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
)

const (
	// Transaction types
	MobileMoney = "mobile-money"
)

type PaymentStore interface {
	CreateCustomer(ctx context.Context, c *Customer) error
	GetUserID(ctx context.Context, cid string) (string, error)
	GetCustomerID(ctx context.Context, uid string) (string, error)
	DeleteCustomer(ctx context.Context, uid string) error
	StorePayment(ctx context.Context, uid, intent string) error
	CreateTransaction(ctx context.Context, t *pb.Transaction) (*pb.Transaction, error)
	GetTransactions(ctx context.Context, uid string) ([]*pb.Transaction, error)
}

type Customer struct {
	UID        string `db:"uid"`
	CustomerID string `db:"customer_id"`
}

type StripeService struct {
	client *client.API
	db     PaymentStore
	logger log.Logger
}

func NewStripeService(db PaymentStore, client *client.API, logger log.Logger, opts ...func(service *StripeService) error) (PaymentService, error) {
	svc := &StripeService{
		client: client,
		db:     db,
		logger: logger,
	}

	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, err
		}
	}

	return svc, nil
}

func (s *StripeService) GetCustomerID(ctx context.Context, uid string) (string, error) {
	return s.db.GetCustomerID(ctx, uid)
}

func (s *StripeService) GetTransactions(ctx context.Context, uid string) ([]*pb.Transaction, error) {
	return s.db.GetTransactions(ctx, uid)
}

func (s *StripeService) SaveCard(ctx context.Context, uid string) (string, error) {
	id, err := s.db.GetCustomerID(ctx, uid)

	if err != nil {
		return "", errors.Wrap(err, "error getting customer id")
	}

	params := &stripe.SetupIntentParams{
		Customer: stripe.String(id),
		Usage:    stripe.String("on_session"),
	}

	intent, err := s.client.SetupIntents.New(params)

	if err != nil {
		return "", errors.Wrap(err, "error creating setup intent")
	}

	return intent.ClientSecret, err
}

func (s *StripeService) CreateTransaction(ctx context.Context, t *pb.Transaction) (*pb.Transaction, error) {
	return s.db.CreateTransaction(ctx, t)
}

func (s *StripeService) CapturePayment(_ context.Context, pi string, amount float64) (string, error) {
	params := &stripe.PaymentIntentCaptureParams{
		AmountToCapture: stripe.Int64(int64(amount * 100)),
	}

	p, err := s.client.PaymentIntents.Capture(pi, params)

	if err != nil {
		return "", err
	}

	return p.ClientSecret, nil
}

func (s *StripeService) GetPaymentIntentSecret(ctx context.Context, req *pb.PaymentRequest) (string, error) {
	id, err := s.db.GetCustomerID(ctx, req.Uid)

	if err != nil {
		return "", err
	}

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(req.Amount * 100)),
		Currency:      stripe.String(req.Currency),
		Customer:      stripe.String(id),
		PaymentMethod: stripe.String(req.CardID),
	}

	if !req.Capture {
		params.CaptureMethod = stripe.String(string(stripe.PaymentIntentCaptureMethodManual))
	}

	pi, err := s.client.PaymentIntents.New(params)

	if err != nil {
		return "", err
	}

	return pi.ClientSecret, nil
}

func (s *StripeService) createCustomer(ctx context.Context, u *userPb.User) (string, error) {
	cus := &Customer{UID: u.Uuid}
	name := fmt.Sprintf("%s %s", u.FirstName, u.LastName)

	params := &stripe.CustomerParams{
		Name:  &name,
		Email: &u.Email,
	}

	c, err := s.client.Customers.New(params)

	if err != nil {
		return "", errors.Wrap(err, "error creating stripe customer")
	}

	cus.CustomerID = c.ID
	err = s.db.CreateCustomer(ctx, cus)

	return c.ID, err
}

func (s *StripeService) deleteCustomer(ctx context.Context, u *userPb.User) error {
	cid, err := s.db.GetCustomerID(ctx, u.Uuid)

	if err != nil {
		return err
	}

	_, err = s.client.Customers.Del(cid, nil)

	if err != nil {
		return err
	}

	return s.db.DeleteCustomer(ctx, u.Uuid)
}

func (s *StripeService) OnUserCreated(ctx context.Context, data *eventpb.EventData) error {
	u := data.GetUser()

	if u == nil {
		return &events.ErrorShouldAck{Err: "error: got nil user"}
	}

	_, err := s.createCustomer(ctx, u)

	return err
}

func (s *StripeService) OnPaymentSucceded(ctx context.Context, data *eventpb.EventData) error {
	intent := data.GetIntent()
	pi, err := s.client.PaymentIntents.Get(intent, nil)

	if err != nil {
		return err
	}

	customerId := pi.Customer.ID
	uid, err := s.db.GetUserID(ctx, customerId)

	if err != nil {
		return err
	}

	return s.db.StorePayment(ctx, intent, uid)
}

func (s *StripeService) OnTransferSucceded(ctx context.Context, data *eventpb.EventData) error {
	tr := data.GetTransfer()
	return s.logger.Log("Transfer", tr)
}
