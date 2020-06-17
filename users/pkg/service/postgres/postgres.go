package postgres

import (
	"context"
	"fmt"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/alexmeli100/remit/users/pkg/service"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

type PostgService struct {
	DB *sqlx.DB
}

func NewPostgService() service.UsersService {
	pass := os.Getenv("POSTGRES_PASSWORD")
	userName := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	connString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", userName, pass, dbName)

	db, err := sqlx.Open("postgres", connString)

	if err != nil {
		log.Fatal("failed to open database connection")
	}

	return &PostgService{DB: db}
}

func (s *PostgService) GetUserByID(ctx context.Context, id int64) (*pb.User, error) {
	u := &pb.User{Id: id}

	err := s.DB.Get(u, "SELECT * FROM users WHERE id=$1 Limit 1", u.Email)

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

func (s *PostgService) UpdatePassword(ctx context.Context, u *pb.User) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.Wrap(err, "bcrypt error")
	}

	_, err = s.DB.Exec("UPDATE users SET password=$1 WHERE email=$2", string(pass), u.Email)

	return err
}

func (s *PostgService) Create(ctx context.Context, u *pb.User) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.Wrap(err, "bcrypt error")
	}

	u.Password = string(pass)

	_, err = s.DB.NamedExec(
		`INSERT INTO users(first_name, last_name, email, address, password, id, confirmed) 
		values(:firstName, :lastName, :email, :address, :password, :id, FALSE) `, u)

	return err
}
