package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	userPb "github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
	"github.com/stripe/stripe-go/v71/customer"
	"github.com/stripe/stripe-go/v71/setupintent"
)

type Customer struct {
	UID        string `db:"uid"`
	CustomerID string `db:"customer_id"`
}

type StripeService struct {
	client *client.API
	db     *sqlx.DB
}

func NewStripeService(db *sqlx.DB, client *client.API, opts ...func(service *StripeService) error) (PaymentService, error) {
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
	id, err := s.getCustomerID(ctx, uid)

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
	id, err := s.getCustomerID(ctx, req.Uid)

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

func (s *StripeService) getCustomerID(ctx context.Context, uid string) (string, error) {
	u := &Customer{UID: uid}

	err := s.db.Get(u, "SELECT * FROM customers WHERE uid=$1", u.UID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.createCustomer(ctx, &userPb.User{Uuid: uid})
		}

		return "", err
	}

	return u.CustomerID, nil
}

func (s *StripeService) createCustomer(_ context.Context, u *userPb.User) (string, error) {
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
	_, err = s.db.NamedExec("INSERT INTO customers(uid, customer_id) values(:customerID, :uID)", cus)

	return c.ID, err
}

func (s *StripeService) OnUserCreated(ctx context.Context, u *userPb.User) error {
	_, err := s.createCustomer(ctx, u)

	return err
}
