package grpc

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/endpoint"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
)

// makeCreateHandler creates the handler logic
func makeCreateHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.CreateEndpoint, decodeCreateRequest, encodeCreateResponse, options...)
}

// decodeCreateRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Create request.
func decodeCreateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CreateRequest)

	return endpoint.CreateRequest{User: req.User}, nil
}

// encodeCreateResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeCreateResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.CreateResponse)

	return &pb.CreateReply{Err: err2str(res.Err)}, nil
}

func (g *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	_, rep, err := g.create.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.CreateReply), nil
}

// makeGetUserByIDHandler creates the handler logic
func makeGetUserByIDHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.GetUserByIDEndpoint, decodeGetUserByIDRequest, encodeGetUserByIDResponse, options...)
}

// decodeGetUserByIDResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetUserByID request.
func decodeGetUserByIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetUserByIDRequest)

	return endpoint.GetUserByIDRequest{Id: req.Id}, nil
}

// encodeGetUserByIDResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGetUserByIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetUserByIDResponse)

	return &pb.GetUserByIDReply{Err: err2str(res.Err), User: res.User}, nil
}

func (g *grpcServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDReply, error) {
	_, rep, err := g.getUserByID.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.GetUserByIDReply), nil
}

// makeGetUserByUUIDHandler creates the handler logic
func makeGetUserByUUIDHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.GetUserByUUIDEndpoint, decodeGetUserByUUIDRequest, encodeGetUserByUUIDResponse, options...)
}

// decodeGetUserByUUIDResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetUserByID request.
func decodeGetUserByUUIDRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetUserByUUIDRequest)

	return endpoint.GetUserByUUIDRequest{UUID: req.UUID}, nil
}

// encodeGetUserByUUIDResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGetUserByUUIDResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetUserByUUIDResponse)

	return &pb.GetUserByUUIDReply{Err: err2str(res.Err), User: res.User}, nil
}

func (g *grpcServer) GetUserByUUID(ctx context.Context, req *pb.GetUserByUUIDRequest) (*pb.GetUserByUUIDReply, error) {
	_, rep, err := g.getUserByUUID.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.GetUserByUUIDReply), nil
}

// makeGetUserByEmailHandler creates the handler logic
func makeGetUserByEmailHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.GetUserByEmailEndpoint, decodeGetUserByEmailRequest, encodeGetUserByEmailResponse, options...)
}

// decodeGetUserByEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain GetUserByEmail request.
func decodeGetUserByEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetUserByEmailRequest)

	return endpoint.GetUserByEmailRequest{Email: req.Email}, nil
}

// encodeGetUserByEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeGetUserByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetUserByEmailResponse)

	return &pb.GetUserByEmailReply{User: res.User, Err: err2str(res.Err)}, nil
}

func (g *grpcServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailReply, error) {
	_, rep, err := g.getUserByEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetUserByEmailReply), nil
}

// makeUpdateEmailHandler creates the handler logic
func makeUpdateEmailHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.UpdateEmailEndpoint, decodeUpdateEmailRequest, encodeUpdateEmailResponse, options...)
}

// decodeUpdateEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain UpdateEmail request.
func decodeUpdateEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UpdateEmailRequest)

	return endpoint.UpdateEmailRequest{User: req.User}, nil
}

// encodeUpdateEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeUpdateEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.UpdateEmailResponse)

	return &pb.UpdateEmailReply{Err: err2str(res.Err)}, nil
}

func (g *grpcServer) UpdateEmail(ctx context.Context, req *pb.UpdateEmailRequest) (*pb.UpdateEmailReply, error) {
	_, rep, err := g.updateEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateEmailReply), nil
}

func makeSetUserProfileHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.SetUserProfileEndpoint, decodeSetUserProfileRequest, endcodeSetUserProfileResponse, options...)
}

func decodeSetUserProfileRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SetUserProfileRequest)

	return endpoint.SetUserProfileRequest{User: req.User}, nil
}

func endcodeSetUserProfileResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.SetUserProfileResponse)

	return &pb.SetUserProfileReply{Err: err2str(res.Err), User: res.User}, nil
}

func (g *grpcServer) SetUserProfile(ctx context.Context, request *pb.SetUserProfileRequest) (*pb.SetUserProfileReply, error) {
	_, rep, err := g.setUserProfile.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.SetUserProfileReply), nil
}

func makeUpdateUserProfileHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.UpdateUserProfileEndpoint, decodeUpdateUserProfileEndpoint, encodeUpdateUserProfileEndpoint, options...)
}

func decodeUpdateUserProfileEndpoint(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UpdateUserProfileRequest)

	return endpoint.UpdateUserProfileRequest{User: req.User}, nil
}

func encodeUpdateUserProfileEndpoint(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.UpdateUserProfileResponse)

	return &pb.UpdateUserProfileReply{User: res.User, Err: err2str(res.Err)}, nil
}

func (g *grpcServer) UpdateUserProfile(ctx context.Context, request *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileReply, error) {
	_, rep, err := g.updateUserProfile.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.UpdateUserProfileReply), nil
}

func makeCreateContactHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.CreateContactEndpoint, decodeCreateContactRequest, encodeCreateContactResponse, options...)
}

func decodeCreateContactRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CreateContactRequest)

	return endpoint.CreateContactRequest{Contact: req.Contact}, nil
}

func encodeCreateContactResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.CreateContactResponse)

	return &pb.CreateContactReply{Err: err2str(res.Err), Contact: res.Contact}, nil
}

func (g *grpcServer) CreateContact(ctx context.Context, req *pb.CreateContactRequest) (*pb.CreateContactReply, error) {
	_, rep, err := g.createContact.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.CreateContactReply), nil
}

func makeGetContactsHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.GetContactsEndpoint, decodeGetContactsRequest, encodeGetContactsResponse, options...)
}

func decodeGetContactsRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetContactsRequest)

	return endpoint.GetContactsRequest{UserId: req.UserID}, nil
}

func encodeGetContactsResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.GetContactsResponse)

	return &pb.GetContactsReply{Contacts: res.Contacts, Err: err2str(res.Err)}, nil
}

func (g *grpcServer) GetContacts(ctx context.Context, req *pb.GetContactsRequest) (*pb.GetContactsReply, error) {
	_, rep, err := g.getContacts.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.GetContactsReply), nil
}

func makeUpdateContactHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.UpdateContactEndpoint, decodeUpdateContactRequest, encodeUpdateContactResponse, options...)
}

func decodeUpdateContactRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UpdateContactRequest)

	return endpoint.UpdateContactRequest{Contact: req.Contact}, nil
}

func encodeUpdateContactResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.UpdateContactResponse)

	return &pb.UpdateContactReply{Contact: res.Contact, Err: err2str(res.Err)}, nil
}

func (g *grpcServer) UpdateContact(ctx context.Context, request *pb.UpdateContactRequest) (*pb.UpdateContactReply, error) {
	_, rep, err := g.updateContact.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.UpdateContactReply), nil
}

func makeDeleteContactHandler(endpoints endpoint.Endpoints, options []grpcTrans.ServerOption) grpcTrans.Handler {
	return grpcTrans.NewServer(endpoints.DeleteContactEndpoint, decodeDeleteContactRequest, encodeDeleteContactResponse, options...)
}

func decodeDeleteContactRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.DeleteContactRequest)

	return endpoint.DeleteContactRequest{Contact: req.Contact}, nil
}

func encodeDeleteContactResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(endpoint.DeleteContactResponse)

	return &pb.DeleteContactReply{Err: err2str(res.Err)}, nil
}

func (g *grpcServer) DeleteContact(ctx context.Context, request *pb.DeleteContactRequest) (*pb.DeleteContactReply, error) {
	_, rep, err := g.deleteContact.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return rep.(*pb.DeleteContactReply), nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
