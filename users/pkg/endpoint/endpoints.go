package endpoint

import (
	service "github.com/alexmeli100/remit/users/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
)

const (
	// list of endpoints
	Create         = "Create"
	GetUserById    = "GetUserByID"
	GetUserByUUID  = "GetUserByUUID"
	GetUserByEmail = "GetUserByEmail"
	UpdateEmail    = "UpdateEmail"
	UpdateStatus   = "UpdateStatus"
	CreateContact  = "CreateContact"
	GetContacts    = "GetContacts"
)

// Endpoints collects all of the endpoints that compose a user service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	CreateEndpoint         endpoint.Endpoint
	GetUserByIDEndpoint    endpoint.Endpoint
	GetUserByUUIDEndpoint  endpoint.Endpoint
	GetUserByEmailEndpoint endpoint.Endpoint
	UpdateEmailEndpoint    endpoint.Endpoint
	UpdateStatusEndpoint   endpoint.Endpoint
	CreateContactEndpoint  endpoint.Endpoint
	GetContactsEndpoint    endpoint.Endpoint
}

// New returns a Endpoints struct that wraps the provided service, and wires in all of the
// expected endpoint middlewares
func New(s service.UsersService, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateEndpoint:         MakeCreateEndpoint(s),
		GetUserByEmailEndpoint: MakeGetUserByEmailEndpoint(s),
		GetUserByIDEndpoint:    MakeGetUserByIDEndpoint(s),
		GetUserByUUIDEndpoint:  MakeGetUserByUUIDEndpoint(s),
		UpdateEmailEndpoint:    MakeUpdateEmailEndpoint(s),
		UpdateStatusEndpoint:   MakeUpdateStatusEndpoint(s),
		CreateContactEndpoint:  MakeCreateContactEndpoint(s),
		GetContactsEndpoint:    MakeGetContactsEndpoint(s),
	}

	for _, m := range mdw[Create] {
		eps.CreateEndpoint = m(eps.CreateEndpoint)
	}
	for _, m := range mdw[GetUserById] {
		eps.GetUserByIDEndpoint = m(eps.GetUserByIDEndpoint)
	}
	for _, m := range mdw[GetUserByUUID] {
		eps.GetUserByUUIDEndpoint = m(eps.GetUserByUUIDEndpoint)
	}
	for _, m := range mdw[GetUserByEmail] {
		eps.GetUserByEmailEndpoint = m(eps.GetUserByEmailEndpoint)
	}
	for _, m := range mdw[UpdateEmail] {
		eps.UpdateEmailEndpoint = m(eps.UpdateEmailEndpoint)
	}
	for _, m := range mdw[UpdateStatus] {
		eps.UpdateStatusEndpoint = m(eps.UpdateStatusEndpoint)
	}

	for _, m := range mdw[CreateContact] {
		eps.CreateContactEndpoint = m(eps.CreateContactEndpoint)
	}

	for _, m := range mdw[GetContacts] {
		eps.GetContactsEndpoint = m(eps.GetContactsEndpoint)
	}

	return eps
}

// returns the list of endpoints
func GetEndpointList() []string {
	return []string{
		Create,
		GetUserByEmail,
		GetUserByUUID,
		GetUserById,
		UpdateEmail,
		UpdateStatus,
	}
}
