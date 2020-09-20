package service

import (
	"context"
	"fmt"
	"github.com/alexmeli100/remit/transfer/pkg/grpc/pb"
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

func (m *MobileTransfer) Transfer(ctx context.Context, r *pb.TransferRequest) *pb.TransferResponse {
	s, ok := m.Services[r.Service]

	if !ok {
		err := fmt.Errorf("unknown service: %s", r.Service)
		return GetTransferResponse(r, err)
	}

	err := s.SendTo(r.Amount, r.RecipientNumber, r.Currency)

	return GetTransferResponse(r, err)
}

func GetTransferResponse(r *pb.TransferRequest, err error) *pb.TransferResponse {
	res := &pb.TransferResponse{
		Amount:          r.Amount,
		RecipientId:     r.RecipientId,
		Currency:        r.Currency,
		Service:         r.Service,
		ReceiveCurrency: r.ReceiveCurrency,
		ExchangeRate:    r.ExchangeRate,
		SendFee:         r.SendFee,
		ReceiveAmount:   r.ReceiveAmount,
		SenderId:        r.SenderId,
		PaymentIntent:   r.PaymentIntent,
		Status:          "Success",
	}

	if err != nil {
		res.Status = "Failed"
		res.FailReason = err.Error()
	}

	return res
}
