package service

import (
	"context"
	"fmt"
)

const (
	// list of mobile money Services
	MTN    = "MTN"
	Orange = "Orange"
)

type MobileTransfer struct {
	Services map[string]SendMoney
}

func NewMobileTransfer(options ...func(*MobileTransfer)) TransferService {
	services := make(map[string]SendMoney)
	m := &MobileTransfer{Services: services}

	for _, option := range options {
		option(m)
	}

	return m
}

func (m *MobileTransfer) Transfer(_ context.Context, r *TransferRequest) *TransferResponse {
	s, ok := m.Services[r.Service]

	if !ok {
		err := fmt.Errorf("unknown service: %s", r.Service)
		return GetTransferResponse(r, err)
	}

	err := s.SendTo(int(r.ReceiveAmount), r.RecipientNumber, r.ReceiveCurrency)

	return GetTransferResponse(r, err)
}

func GetTransferResponse(r *TransferRequest, err error) *TransferResponse {
	res := &TransferResponse{
		Amount:          r.Amount,
		RecipientId:     r.RecipientId,
		Currency:        r.Currency,
		Service:         r.Service,
		ReceiveCurrency: r.ReceiveCurrency,
		ExchangeRate:    r.ExchangeRate,
		SendFee:         r.SendFee,
		ReceiveAmount:   r.ReceiveAmount,
		SenderId:        r.SenderId,
		Status:          "Success",
	}

	if err != nil {
		res.Status = "Failed"
		res.FailReason = err.Error()
	}

	return res
}
