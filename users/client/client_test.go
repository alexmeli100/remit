package client

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/alexmeli100/remit/users/pkg/service"
	grpcTrans "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
)

var client service.UsersService

func TestMain(m *testing.M) {
	addr := "10.102.208.155:32701"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	opts := make(map[string][]grpcTrans.ClientOption)

	client = NewGRPCClient(conn, opts)

	code := m.Run()
	conn.Close()
	os.Exit(code)
}

func TestCreateEndpoint(t *testing.T) {
	u := &pb.User{
		FirstName: "Alex",
		LastName:  "Meli",
		Email:     "alexmeli100@gmail.com",
		Uuid:      "lkeln;qoi;4",
		Country:   "Canada",
	}

	if err := client.Create(context.Background(), u); err != nil {
		t.Error(err)
	}
}
