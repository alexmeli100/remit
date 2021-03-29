package service

import (
	"context"
)

type SendMoney interface {
	SendTo(amount int, recipient, currency string) error
}

type TransferRequest struct {
	Amount          float64 `json:"amount,omitempty"`
	RecipientId     int64   `json:"recipientId,omitempty"`
	RecipientNumber string  `json:"recipientNumber,omitempty"`
	Currency        string  `json:"currency,omitempty"`
	Service         string  `json:"service,omitempty"`
	ReceiveCurrency string  `json:"receiveCurrency,omitempty"`
	ExchangeRate    float64 `json:"exchangeRate,omitempty"`
	SendFee         float64 `json:"sendFee,omitempty"`
	ReceiveAmount   int64   `json:"receiveAmount,omitempty"`
	SenderId        string  `json:"senderId,omitempty"`
}

type TransferResponse struct {
	Amount          float64 `json:"amount,omitempty"`
	RecipientId     int64   `json:"recipientId,omitempty"`
	RecipientNumber string  `json:"recipientNumber,omitempty"`
	Currency        string  `json:"currency,omitempty"`
	Service         string  `json:"service,omitempty"`
	ReceiveCurrency string  `json:"receiveCurrency,omitempty"`
	ExchangeRate    float64 `json:"exchangeRate,omitempty"`
	SendFee         float64 `json:"sendFee,omitempty"`
	ReceiveAmount   int64   `json:"receiveAmount,omitempty"`
	SenderId        string  `json:"senderId,omitempty"`
	Status          string  `json:"status,omitempty"`
	FailReason      string  `json:"failReason,omitempty"`
}

// TransferService describes the service.
type TransferService interface {
	Transfer(ctx context.Context, request *TransferRequest) *TransferResponse
}

//New returns a TransferService with all of the expected middleware wired in.
func New(svc TransferService, middleware []Middleware) TransferService {

	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
