package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

const (
	MTN    = "MTN"
	Orange = "Orange"
)

const (
	StatusSuccess   = "SUCCESS"
	StatusPending   = "PENDING"
	StatusFailed    = "FAILED"
	StatusExpired   = "EXPIRED"
	StatusInitiated = "INITIATED"
	StatusUnknown   = "UNKNOWN"
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

func (m *MobileTransfer) Transfer(_ context.Context, r *TransferRequest) (*TransferResponse, error) {
	s, ok := m.Services[r.Service]

	if !ok {
		err := fmt.Errorf("unknown service: %s", r.Service)
		return nil, err
	}

	res, err := s.SendTo(r)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to transfer money")
	}

	return res, nil
}
