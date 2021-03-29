package service

import (
	"context"
)

type SendMoney interface {
	SendTo(req *TransferRequest) (*TransferResponse, error)
}

type TransferRequest struct {
	Amount          float64 `json:"amount,omitempty"`
	RecipientId     int64   `json:"recipientId,omitempty"`
	RecipientNumber string  `json:"recipientNumber,omitempty"`
	Currency        string  `json:"currency,omitempty"`
	Service         string  `json:"service,omitempty"`
	ReceiveCurrency string  `json:"receiveCurrency,omitempty"`
	ReceiveAmount   int64   `json:"receiveAmount,omitempty"`
	SenderId        string  `json:"senderId,omitempty"`
	OrderId         string  `json:"orderId,omitempty"`
}

type TransferResponse struct {
	Amount          float64 `json:"amount,omitempty"`
	RecipientId     int64   `json:"recipientId,omitempty"`
	RecipientNumber string  `json:"recipientNumber,omitempty"`
	Currency        string  `json:"currency,omitempty"`
	Service         string  `json:"service,omitempty"`
	ReceiveCurrency string  `json:"receiveCurrency,omitempty"`
	ReceiveAmount   int64   `json:"receiveAmount,omitempty"`
	SenderId        string  `json:"senderId,omitempty"`
	Status          string  `json:"status,omitempty"`
	Token           string  `json:"token,omitempty"`
	RemoteId        string  `json:"remoteId,omitempty"`
	OrderId         string  `json:"orderId,omitempty"`
	FailReason      string  `json:"failReason,omitempty"`
}

// TransferService describes the service.
type TransferService interface {
	Transfer(ctx context.Context, request *TransferRequest) (*TransferResponse, error)
}
