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

func (m *MobileTransfer) Transfer(ctx context.Context, r *pb.TransferRequest) error {
	s, ok := m.Services[r.Service]

	if !ok {
		return fmt.Errorf("unknown service: %s", r.Service)
	}

	return s.SendTo(r.Amount, r.RecipientNumber, r.Currency)
}
