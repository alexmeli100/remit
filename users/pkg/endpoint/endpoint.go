package endpoint

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/alexmeli100/remit/users/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

// CreateRequest collects the request parameters for the Create method.
type CreateRequest struct {
	User *pb.User `json:"user"`
}

// CreateResponse collects the response parameters for the Create method.
type CreateResponse struct {
	Err error `json:"err"`
}

// MakeCreateEndpoint returns an endpoint that invokes Create on the service.
func MakeCreateEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		err := s.Create(ctx, req.User)

		return CreateResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r CreateResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// GetUserByIDRequest collects the request parameters for the GetUserByID method.
type GetUserByIDRequest struct {
	Id int64 `json:"id"`
}

// GetUserByIDResponse collects the response parameters for the GetUserByID method.
type GetUserByIDResponse struct {
	User *pb.User `json:"user"`
	Err  error    `json:"err"`
}

// MakeGetUserByIDEndpoint returns an endpoint that invokes GetUserByID on the service.
func MakeGetUserByIDEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserByIDRequest)
		user, err := s.GetUserByID(ctx, req.Id)

		return GetUserByIDResponse{
			Err:  err,
			User: user,
		}, nil
	}
}

// Failed implements Failer.
func (r GetUserByIDResponse) Failed() error {
	return r.Err
}

// GetUserByIDRequest collects the request parameters for the GetUserByID method.
type GetUserByUUIDRequest struct {
	UUID string `json:"uuid"`
}

// GetUserByIDResponse collects the response parameters for the GetUserByID method.
type GetUserByUUIDResponse struct {
	User *pb.User `json:"user"`
	Err  error    `json:"err"`
}

// MakeGetUserByIDEndpoint returns an endpoint that invokes GetUserByID on the service.
func MakeGetUserByUUIDEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserByUUIDRequest)
		user, err := s.GetUserByUUID(ctx, req.UUID)

		return GetUserByIDResponse{
			Err:  err,
			User: user,
		}, nil
	}
}

// Failed implements Failer.
func (r GetUserByUUIDResponse) Failed() error {
	return r.Err
}

// GetUserByEmailRequest collects the request parameters for the GetUserByEmail method.
type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

// GetUserByEmailResponse collects the response parameters for the GetUserByEmail method.
type GetUserByEmailResponse struct {
	User *pb.User `json:"user"`
	Err  error    `json:"err"`
}

// MakeGetUserByEmailEndpoint returns an endpoint that invokes GetUserByEmail on the service.
func MakeGetUserByEmailEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserByEmailRequest)
		user, err := s.GetUserByEmail(ctx, req.Email)
		res := GetUserByEmailResponse{Err: err, User: user}

		return res, nil
	}
}

// Failed implements Failer.
func (r GetUserByEmailResponse) Failed() error {
	return r.Err
}

// UpdateEmailRequest collects the request parameters for the UpdateEmail method.
type UpdateEmailRequest struct {
	User *pb.User `json:"user"`
}

// UpdateEmailResponse collects the response parameters for the UpdateEmail method.
type UpdateEmailResponse struct {
	Err error `json:"err"`
}

// MakeUpdateEmailEndpoint returns an endpoint that invokes UpdateEmail on the service.
func MakeUpdateEmailEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateEmailRequest)
		err := s.UpdateEmail(ctx, req.User)

		return UpdateEmailResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r UpdateEmailResponse) Failed() error {
	return r.Err
}

// UpdateStatusRequest collects the request parameters for the UpdateStatus method.
type UpdateStatusRequest struct {
	User *pb.User `json:"user"`
}

// UpdateStatusResponse collects the response parameters for the UpdateStatus method.
type UpdateStatusResponse struct {
	Err error `json:"err"`
}

// MakeUpdateStatusEndpoint returns an endpoint that invokes UpdateStatus on the service.
func MakeUpdateStatusEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateStatusRequest)
		err := s.UpdateStatus(ctx, req.User)

		return UpdateStatusResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r UpdateStatusResponse) Failed() error {
	return r.Err
}

// Create implements Service. Primarily useful in a client.
func (e Endpoints) Create(ctx context.Context, user *pb.User) error {
	request := CreateRequest{User: user}
	response, err := e.CreateEndpoint(ctx, request)

	if err != nil {
		return err
	}

	return response.(CreateResponse).Err
}

// GetUserByID implements Service. Primarily useful in a client.
func (e Endpoints) GetUserByID(ctx context.Context, id int64) (*pb.User, error) {
	request := GetUserByIDRequest{Id: id}
	response, err := e.GetUserByIDEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	r := response.(GetUserByIDResponse)

	return r.User, r.Err
}

func (e Endpoints) GetUserByUUID(ctx context.Context, uuid string) (*pb.User, error) {
	request := GetUserByUUIDRequest{UUID: uuid}
	response, err := e.GetUserByUUIDEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	r := response.(GetUserByUUIDResponse)

	return r.User, r.Err
}

// GetUserByEmail implements Service. Primarily useful in a client.
func (e Endpoints) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	request := GetUserByEmailRequest{Email: email}
	response, err := e.GetUserByEmailEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	r := response.(GetUserByEmailResponse)

	return r.User, r.Err
}

// UpdateEmail implements Service. Primarily useful in a client.
func (e Endpoints) UpdateEmail(ctx context.Context, user *pb.User) (e0 error) {
	request := UpdateEmailRequest{User: user}
	response, err := e.UpdateEmailEndpoint(ctx, request)
	if err != nil {
		return
	}

	return response.(UpdateEmailResponse).Err
}

// UpdateStatus implements Service. Primarily useful in a client.
func (e Endpoints) UpdateStatus(ctx context.Context, user *pb.User) (e0 error) {
	request := UpdateStatusRequest{User: user}
	response, err := e.UpdateStatusEndpoint(ctx, request)
	if err != nil {
		return
	}

	return response.(UpdateStatusResponse).Err
}
