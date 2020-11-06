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
	User *pb.User `json:"user"`
	Err  error    `json:"err"`
}

// MakeCreateEndpoint returns an endpoint that invokes Create on the service.
func MakeCreateEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		u, err := s.Create(ctx, req.User)

		return CreateResponse{Err: err, User: u}, nil
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

		return GetUserByUUIDResponse{Err: err, User: user}, nil
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

type CreateContactRequest struct {
	Contact *pb.Contact `json:"contact"`
}

type CreateContactResponse struct {
	Contact *pb.Contact `json:"contact"`
	Err     error       `json:"err"`
}

func MakeCreateContactEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateContactRequest)
		c, err := s.CreateContact(ctx, req.Contact)

		return CreateContactResponse{Err: err, Contact: c}, nil
	}
}

func (r CreateContactResponse) Failed() error {
	return r.Err
}

type UpdateContactRequest struct {
	Contact *pb.Contact `json:"contact"`
}

type UpdateContactResponse struct {
	Contact *pb.Contact `json:"contact"`
	Err     error       `json:"err"`
}

func MakeUpdateContactEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateContactRequest)
		c, err := s.UpdateContact(ctx, req.Contact)

		return UpdateContactResponse{Err: err, Contact: c}, nil
	}
}

func (r UpdateContactResponse) Failed() error {
	return r.Err
}

type SetUserProfileRequest struct {
	User *pb.User `json:"user"`
}

type SetUserProfileResponse struct {
	User *pb.User `json:"user"`
	Err  error    `json:"err"`
}

func MakeSetUserProfileResponse(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SetUserProfileRequest)
		u, err := s.SetUserProfile(ctx, req.User)

		return SetUserProfileResponse{User: u, Err: err}, nil
	}
}

func (r SetUserProfileResponse) Failed() error {
	return r.Err
}

type UpdateUserProfileRequest struct {
	User *pb.User `json:"user"`
}

type UpdateUserProfileResponse struct {
	User *pb.User `json:"user"`
	Err  error    `json:"err"`
}

func MakeUpdateUserProfileEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateUserProfileRequest)
		u, err := s.UpdateUserProfile(ctx, req.User)

		return UpdateUserProfileResponse{User: u, Err: err}, nil
	}
}

type GetContactsRequest struct {
	UserId int64 `json:"userId"`
}

type GetContactsResponse struct {
	Contacts []*pb.Contact `json:"contacts"`
	Err      error         `json:"err"`
}

func MakeGetContactsEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetContactsRequest)
		c, err := s.GetContacts(ctx, req.UserId)

		return GetContactsResponse{Err: err, Contacts: c}, nil
	}
}

type DeleteContactRequest struct {
	Contact *pb.Contact `json:"contact"`
}

type DeleteContactResponse struct {
	Err error `json:"err"`
}

func MakeDeleteContactEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteContactRequest)
		err := s.DeleteContact(ctx, req.Contact)

		return DeleteContactResponse{Err: err}, nil
	}
}

// Create implements Service. Primarily useful in a client.
func (e Endpoints) Create(ctx context.Context, user *pb.User) (*pb.User, error) {
	request := CreateRequest{User: user}
	response, err := e.CreateEndpoint(ctx, request)

	if err != nil {
		return nil, err
	}

	res := response.(CreateResponse)

	return res.User, res.Err
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
func (e Endpoints) UpdateEmail(ctx context.Context, user *pb.User) error {
	request := UpdateEmailRequest{User: user}
	response, err := e.UpdateEmailEndpoint(ctx, request)
	if err != nil {
		return err
	}

	return response.(UpdateEmailResponse).Err
}

func (e Endpoints) CreateContact(ctx context.Context, contact *pb.Contact) (*pb.Contact, error) {
	request := CreateContactRequest{Contact: contact}
	response, err := e.CreateContactEndpoint(ctx, request)

	if err != nil {
		return nil, err
	}

	res := response.(CreateContactResponse)

	return res.Contact, res.Err
}

func (e Endpoints) GetContacts(ctx context.Context, userId int64) ([]*pb.Contact, error) {
	request := GetContactsRequest{UserId: userId}
	response, err := e.GetContactsEndpoint(ctx, request)

	if err != nil {
		return nil, err
	}

	res := response.(GetContactsResponse)
	return res.Contacts, res.Err
}

func (e Endpoints) UpdateContact(ctx context.Context, contact *pb.Contact) (*pb.Contact, error) {
	request := UpdateContactRequest{Contact: contact}
	response, err := e.UpdateContactEndpoint(ctx, request)

	if err != nil {
		return nil, err
	}

	res := response.(UpdateContactResponse)
	return res.Contact, res.Err
}

func (e Endpoints) DeleteContact(ctx context.Context, contact *pb.Contact) error {
	request := DeleteContactRequest{Contact: contact}
	response, err := e.DeleteContactEndpoint(ctx, request)

	if err != nil {
		return err
	}

	res := response.(DeleteContactResponse)
	return res.Err
}

func (e Endpoints) SetUserProfile(ctx context.Context, user *pb.User) (*pb.User, error) {
	request := SetUserProfileRequest{User: user}
	response, err := e.SetUserProfileEndpoint(ctx, request)

	if err != nil {
		return nil, err
	}

	res := response.(SetUserProfileResponse)
	return res.User, res.Err
}

func (e Endpoints) UpdateUserProfile(ctx context.Context, user *pb.User) (*pb.User, error) {
	request := UpdateUserProfileRequest{User: user}
	response, err := e.UpdateUserProfileEndpoint(ctx, request)

	if err != nil {
		return nil, err
	}

	res := response.(UpdateUserProfileResponse)
	return res.User, res.Err
}
