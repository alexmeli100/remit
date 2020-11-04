package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const (
	CreateCustomerQuery    = "INSERT INTO customers(uid, customer_id) values(:uid, :customer_id)"
	GetUserIDQuery         = "SELECT * FROM customers WHERE customer_id=$1"
	GetCustomerIDQuery     = "SELECT * FROM customers WHERE uid=$1"
	StorePaymentQuery      = "INSERT INTO payments(uid, intent) values($1, $2)"
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
	db *sqlx.DB
}

func NewPostgresDB(db *sqlx.DB) PaymentStore {
	return &PostgresDB{db}
}

func (p *PostgresDB) CreateCustomer(ctx context.Context, c *Customer) error {
	_, err := p.db.NamedExec(CreateCustomerQuery, c)

	return err
}

func (p *PostgresDB) GetUserID(ctx context.Context, cid string) (string, error) {
	var c Customer
	err := p.db.Get(&c, GetUserIDQuery, cid)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoUser
	} else if err != nil {
		return "", err
	}

	return c.UID, nil
}

func (p *PostgresDB) DeleteCustomer(ctx context.Context, uid string) error {
	_, err := p.db.Exec(DeleteCustomerQuery, uid)

	return err
}

func (p *PostgresDB) GetCustomerID(ctx context.Context, uid string) (string, error) {
	var c Customer
	err := p.db.Get(&c, GetCustomerIDQuery, uid)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoUser
	} else if err != nil {
		return "", err
	}

	return c.CustomerID, nil
}

func (p *PostgresDB) StorePayment(ctx context.Context, uid, intent string) error {
	_, err := p.db.Exec(StorePaymentQuery, uid, intent)

	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) CreateTransaction(ctx context.Context, t *pb.Transaction) (*pb.Transaction, error) {
	var lastInsertId int

	err := p.db.QueryRowx(CreateTransactionQuery,
		t.RecipientId, t.UserId, time.Now(), t.AmountReceived, t.AmountSent, t.TransactionFee,
		t.TransactionType, t.SendCurrency, t.ReceiveCurrency, t.ExchangeRate, t.PaymentIntent,
	).Scan(&lastInsertId)

	if err != nil {
		return nil, err
	}

	t.Id = int64(lastInsertId)
	return t, nil
}

func (p *PostgresDB) GetTransactions(ctx context.Context, uid string) ([]*pb.Transaction, error) {
	var trs []*pb.Transaction

	err := p.db.Select(trs, GetTransactionsQuery, uid)

	return trs, err
}
