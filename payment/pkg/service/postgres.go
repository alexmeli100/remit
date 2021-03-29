package service

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const (
	GetTransactionsQuery   = "SELECT * FROM transactions WHERE user_id=$1"
	CreateTransactionQuery = `
		INSERT INTO transactions(
			recipient_id, user_id, created_at, amount_received, amount_sent, transaction_fee, 
			transaction_type, send_currency, receive_currency, exchange_rate, payment_intent)
		values(
			$1, $2, $3, $4, $5, $6,$7, $8, $9, $10, $11)
		RETURNING id
		`
	DeleteCustomerQuery = "DELETE FROM customers WHERE uid=$1"
)

var ErrNoUser = errors.New("user not found")

type PostgresDB struct {
	DB *sqlx.DB
}

func NewPostgresDB(db *sqlx.DB) PaymentStore {
	return &PostgresDB{db}
}

func (p *PostgresDB) CreateTransaction(ctx context.Context, t *Transaction) (*Transaction, error) {
	var lastInsertId int

	err := p.DB.QueryRowxContext(ctx, CreateTransactionQuery,
		t.RecipientId, t.UserId, time.Now(), t.AmountReceived, t.AmountSent, t.TransactionFee,
		t.TransactionType, t.SendCurrency, t.ReceiveCurrency, t.ExchangeRate, t.PaymentIntent,
	).Scan(&lastInsertId)

	if err != nil {
		return nil, err
	}

	t.Id = int64(lastInsertId)
	return t, nil
}

func (p *PostgresDB) GetTransactions(ctx context.Context, uid string) ([]*Transaction, error) {
	var trs []*Transaction

	err := p.DB.SelectContext(ctx, trs, GetTransactionsQuery, uid)

	return trs, err
}
