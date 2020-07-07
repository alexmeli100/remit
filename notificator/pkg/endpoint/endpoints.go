package endpoint

import (
	service "github.com/alexmeli100/remit/notificator/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
)

const (
	SendConfirmEmail       = "SendConfirmEmail"
	SendPasswordResetEmail = "SendPasswordResetEmail"
	SendWelcomeEmail       = "SendWelcomeEmail"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	SendConfirmEmailEndpoint       endpoint.Endpoint
	SendPasswordResetEmailEndpoint endpoint.Endpoint
	SendWelcomeEmailEndpoint       endpoint.Endpoint
}

// New returns a Endpoints struct that wraps the provided service, and wires in all of the
// expected endpoint middlewares
func New(s service.NotificatorService, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		SendConfirmEmailEndpoint:       MakeSendConfirmEmailEndpoint(s),
		SendPasswordResetEmailEndpoint: MakeSendPasswordResetEmailEndpoint(s),
		SendWelcomeEmailEndpoint:       MakeSendWelcomeEmailEndpoint(s),
	}

	for _, m := range mdw[SendConfirmEmail] {
		eps.SendConfirmEmailEndpoint = m(eps.SendConfirmEmailEndpoint)
	}
	for _, m := range mdw[SendPasswordResetEmail] {
		eps.SendPasswordResetEmailEndpoint = m(eps.SendPasswordResetEmailEndpoint)
	}
	for _, m := range mdw[SendWelcomeEmail] {
		eps.SendWelcomeEmailEndpoint = m(eps.SendWelcomeEmailEndpoint)
	}
	return eps
}

func GetEndpointList() []string {
	return []string{
		SendConfirmEmail,
		SendPasswordResetEmail,
		SendWelcomeEmail,
	}
}
