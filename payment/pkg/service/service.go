package service

import (
	"context"
	"time"
)

type Transaction struct {
	Id              int64      `json:"id,omitempty" db:"id"`
	RecipientId     int64      `json:"recipientId,omitempty" db:"recipient_id"`
	UserId          string     `json:"userId,omitempty" db:"user_id"`
	CreatedAt       *time.Time `json:"createdAt,omitempty" db:"created_at"`
	AmountReceived  float64    `json:"amountReceived,omitempty" db:"amount_received"`
	AmountSent      float64    `json:"amountSent,omitempty" db:"amount_sent"`
	TransactionFee  float64    `json:"transactionFee,omitempty" db:"transaction_fee"`
	TransactionType string     `json:"transactionType,omitempty" db:"transaction_type"`
	SendCurrency    string     `json:"sendCurrency,omitempty" db:"send_currency"`
	ReceiveCurrency string     `json:"receiveCurrency,omitempty" db:"receive_currency"`
	ExchangeRate    float64    `json:"exchangeRate,omitempty" db:"exchange_rate"`
	PaymentIntent   string     `json:"paymentIntent,omitempty" db:"payment_intent"`
}

type PaymentStore interface {
	CreateTransaction(ctx context.Context, tr *Transaction) (*Transaction, error)
	GetTransactions(ctx context.Context, uid string) ([]*Transaction, error)
}

// PaymentService describes the service.
type PaymentService interface {
	CreateTransaction(ctx context.Context, tr *Transaction) (*Transaction, error)
	GetTransactions(ctx context.Context, uid string) ([]*Transaction, error)
}
