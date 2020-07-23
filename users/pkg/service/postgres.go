package service

import (
	"context"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgService struct {
	DB *sqlx.DB
}

func NewPostgService(db *sqlx.DB) UsersService {
	return &PostgService{DB: db}
}

func (s *PostgService) GetUserByID(ctx context.Context, id int64) (*pb.User, error) {
	u := &pb.User{Id: id}

	err := s.DB.Get(u, "SELECT * FROM users WHERE id=$1 Limit 1", u.Id)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *PostgService) GetUserByUUID(ctx context.Context, uuid string) (*pb.User, error) {
	u := &pb.User{Uuid: uuid}

	err := s.DB.Get(u, "SELECT * FROM users WHERE uuid=$1 Limit 1", u.Uuid)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *PostgService) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	u := &pb.User{Email: email}

	err := s.DB.Get(u, "SELECT * FROM users WHERE email=$1 Limit 1", u.Email)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *PostgService) UpdateEmail(ctx context.Context, u *pb.User) error {
	_, err := s.DB.Exec("UPDATE users SET email=$1 WHERE id=$2", u.Email, u.Id)

	return err
}

func (s *PostgService) UpdateStatus(ctx context.Context, u *pb.User) error {
	_, err := s.DB.Exec("UPDATE users SET confirmed=TRUE WHERE id=$2", u.Id)

	return err
}

func (s *PostgService) Create(ctx context.Context, u *pb.User) error {
	_, err := s.DB.Exec(
		`INSERT INTO users(first_name, last_name, email, country, uuid, confirmed) 
		values($1, $2, $3, $4, $5, FALSE) `, u.FirstName, u.LastName, u.Email, u.Country, u.Uuid)

	return err
}
