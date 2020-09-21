package service

import (
	"context"
	"fmt"
	eventpb "github.com/alexmeli100/remit/events/pb"
	"github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	userPb "github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
	"github.com/stripe/stripe-go/v71/customer"
	"github.com/stripe/stripe-go/v71/setupintent"
	"time"
)

const (
	// Transaction types
	MobileMoney = "mobile-money"
)

type PaymentStore interface {
	CreateCustomer(ctx context.Context, c *Customer) error
	GetUserID(ctx context.Context, cid string) (string, error)
	GetCustomerID(ctx context.Context, uid string) (string, error)
	StorePayment(ctx context.Context, uid, intent string) error
	CreateTransaction(ctx context.Context, t *pb.Transaction) error
	GetTransactions(ctx context.Context, uid string) ([]*pb.Transaction, error)
}

type Customer struct {
	UID        string `db:"uid"`
	CustomerID string `db:"customer_id"`
}

type StripeService struct {
	client *client.API
	db     PaymentStore
}

func NewStripeService(db PaymentStore, client *client.API, opts ...func(service *StripeService) error) (PaymentService, error) {
	svc := &StripeService{
		client: client,
		db:     db,
	}

	for _, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, err
		}
	}

	return svc, nil
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

	intent, err := setupintent.New(params)

	if err != nil {
		return "", errors.Wrap(err, "error creating setup intent")
	}

	return intent.ClientSecret, err
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

	c, err := customer.New(params)

	if err != nil {
		return "", errors.Wrap(err, "error creating stripe customer")
	}

	cus.UID = c.ID
	err = s.db.CreateCustomer(ctx, cus)

	return c.ID, err
}

func (s *StripeService) OnUserCreated(ctx context.Context, data *eventpb.EventData) error {
	u := data.GetUser()
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
	now := time.Now()

	order := &pb.Transaction{
		RecipientId:     tr.RecipientId,
		UserId:          tr.SenderId,
		CreatedAt:       &now,
		AmountReceived:  tr.ReceiveAmount,
		AmountSent:      tr.Amount,
		TransactionFee:  tr.SendFee,
		TransactionType: MobileMoney,
		SendCurrency:    tr.Currency,
		ReceiveCurrency: tr.ReceiveCurrency,
		ExchangeRate:    tr.ExchangeRate,
		PaymentIntent:   tr.PaymentIntent,
	}

	return s.db.CreateTransaction(ctx, order)
}
