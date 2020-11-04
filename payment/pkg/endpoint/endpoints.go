package endpoint

import (
	"github.com/alexmeli100/remit/payment/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

const (
	// list of endpoints
	GetPaymentIntentSecret = "GetPaymentIntentSecret"
	SaveCard               = "SaveCard"
	CapturePayment         = "CapturePayment"
	GetCustomerID          = "GetCustomerID"
	CreateTransaction      = "CreateTransaction"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	SaveCardEndpoint               endpoint.Endpoint
	GetPaymentIntentSecretEndpoint endpoint.Endpoint
	CapturePaymentEndpoint         endpoint.Endpoint
	GetCustomerIDEndpoint          endpoint.Endpoint
	CreateTransactionEndpoint      endpoint.Endpoint
}

// New returns a Endpoints struct that wraps the provided service, and wires in all of the
// expected endpoint middlewares
func New(s service.PaymentService, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GetPaymentIntentSecretEndpoint: MakeGetPaymentIntentSecretEndpoint(s),
		SaveCardEndpoint:               MakeSaveCardEndpoint(s),
		CapturePaymentEndpoint:         MakeCapturePaymentEndpoint(s),
	}

	for _, m := range mdw[SaveCard] {
		eps.SaveCardEndpoint = m(eps.SaveCardEndpoint)
	}
	for _, m := range mdw[GetPaymentIntentSecret] {
		eps.GetPaymentIntentSecretEndpoint = m(eps.GetPaymentIntentSecretEndpoint)
	}
	for _, m := range mdw[CapturePayment] {
		eps.CapturePaymentEndpoint = m(eps.CapturePaymentEndpoint)
	}

	return eps
}

func GetEndpointList() []string {
	return []string{
		SaveCard,
		GetPaymentIntentSecret,
		CapturePayment,
		GetCustomerID,
		CreateTransaction,
	}
}
