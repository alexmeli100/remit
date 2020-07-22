package service

import "context"

const (
	// list of operators
	MTN    = "mtn"
	Orange = "orange"
)

type SendOptions struct {
	operator  string
	recipient string
	amount    string
}

// TransferService describes the service.
type TransferService interface {
	Send(ctx context.Context, req *SendOptions) error
}

type basicTransferService struct{}

func (b *basicTransferService) Send(ctx context.Context, req *SendOptions) (e0 error) {
	// TODO implement the business logic of Send
	return e0
}

// NewBasicTransferService returns a naive, stateless implementation of TransferService.
func NewBasicTransferService() TransferService {
	return &basicTransferService{}
}

// New returns a TransferService with all of the expected middleware wired in.
func New(middleware []Middleware) TransferService {
	var svc TransferService = NewBasicTransferService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
