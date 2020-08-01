package endpoint

import (
	"context"

	service "github.com/alexmeli100/remit/notificator/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
)

// SendConfirmEmailRequest collects the request parameters for the SendConfirmEmail method.
type SendConfirmEmailRequest struct {
	Addr string `json:"addr"`
	Link string `json:"link"`
	Name string `json:"name"`
}

// SendConfirmEmailResponse collects the response parameters for the SendConfirmEmail method.
type SendConfirmEmailResponse struct {
	Err error `json:"err"`
}

// MakeSendConfirmEmailEndpoint returns an endpoint that invokes SendConfirmEmail on the service.
func MakeSendConfirmEmailEndpoint(s service.NotificatorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendConfirmEmailRequest)
		err := s.SendConfirmEmail(ctx, req.Name, req.Addr, req.Link)
		return SendConfirmEmailResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r SendConfirmEmailResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// SendConfirmEmail implements Service. Primarily useful in a client.
func (e Endpoints) SendConfirmEmail(ctx context.Context, name, addr string, link string) error {
	request := SendConfirmEmailRequest{
		Addr: addr,
		Name: name,
		Link: link,
	}

	response, err := e.SendConfirmEmailEndpoint(ctx, request)

	if err != nil {
		return err
	}
	return response.(SendConfirmEmailResponse).Err
}

// SendPasswordResetEmailRequest collects the request parameters for the SendPasswordResetEmail method.
type SendPasswordResetEmailRequest struct {
	Addr string `json:"addr"`
	Link string `json:"link"`
}

// SendPasswordResetEmailResponse collects the response parameters for the SendPasswordResetEmail method.
type SendPasswordResetEmailResponse struct {
	Err error `json:"err"`
}

// MakeSendPasswordResetEmailEndpoint returns an endpoint that invokes SendPasswordResetEmail on the service.
func MakeSendPasswordResetEmailEndpoint(s service.NotificatorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendPasswordResetEmailRequest)
		err := s.SendPasswordResetEmail(ctx, req.Addr, req.Link)
		return SendPasswordResetEmailResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r SendPasswordResetEmailResponse) Failed() error {
	return r.Err
}

// SendWelcomeEmailRequest collects the request parameters for the SendWelcomeEmail method.
type SendWelcomeEmailRequest struct {
	Addr string `json:"addr"`
	Name string `json:"name"`
}

// SendWelcomeEmailResponse collects the response parameters for the SendWelcomeEmail method.
type SendWelcomeEmailResponse struct {
	Err error `json:"err"`
}

// MakeSendWelcomeEmailEndpoint returns an endpoint that invokes SendWelcomeEmail on the service.
func MakeSendWelcomeEmailEndpoint(s service.NotificatorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendWelcomeEmailRequest)
		err := s.SendWelcomeEmail(ctx, req.Name, req.Addr)

		return SendWelcomeEmailResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r SendWelcomeEmailResponse) Failed() error {
	return r.Err
}

// SendPasswordResetEmail implements Service. Primarily useful in a client.
func (e Endpoints) SendPasswordResetEmail(ctx context.Context, addr string, link string) error {
	request := SendPasswordResetEmailRequest{
		Addr: addr,
		Link: link,
	}
	response, err := e.SendPasswordResetEmailEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return response.(SendPasswordResetEmailResponse).Err
}

// SendWelcomeEmail implements Service. Primarily useful in a client.
func (e Endpoints) SendWelcomeEmail(ctx context.Context, name, addr string) error {
	request := SendWelcomeEmailRequest{
		Addr: addr,
		Name: name,
	}
	response, err := e.SendWelcomeEmailEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return response.(SendWelcomeEmailResponse).Err
}
