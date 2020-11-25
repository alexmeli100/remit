package service

import (
	"context"
	"github.com/alexmeli100/remit/payment/pkg/grpc/pb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

var pg PostgresDB

func clearTransactions() {
	if _, err := pg.DB.Exec("DELETE FROM transactions"); err != nil {
		log.Fatal(err)
	}
}

func TestPostGres(t *testing.T) {
	db, err := openConnection()

	if err != nil {
		log.Fatal(err)
	}

	//defer db.Close()
	pg = PostgresDB{DB: db}

	t.Run("testPostgres", func(t *testing.T) {
		t.Parallel()
		TestCreateTransaction(t)
	})
}

func TestCreateTransaction(t *testing.T) {
	clearTransactions()
	now := time.Now()

	tr := &pb.Transaction{
		RecipientId:     2,
		UserId:          uuid.New().String(),
		CreatedAt:       &now,
		AmountReceived:  0,
		AmountSent:      0,
		TransactionFee:  0,
		TransactionType: "mobile",
		SendCurrency:    "CAD",
		ReceiveCurrency: "XAF",
		ExchangeRate:    0,
		PaymentIntent:   uuid.New().String(),
	}

	newTr, err := pg.CreateTransaction(context.Background(), tr)

	if assert.NoError(t, err) {
		assert.NotNil(t, newTr)
	}
}
