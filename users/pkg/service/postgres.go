package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const (
	UpdateEmailQuery = "UPDATE users SET email=$1 WHERE id=$2"

	getUserByQuery = `SELECT
						    users.id,
							users.uuid,
						    users.first_name,
						    users.middle_name,
						    users.last_name,
						    users.email,
						    users.confirmed,
						    users.created_at,
						    users.country,
							profile.birth_date,
							profile.gender,
							profile.occupation,
							user_address.country,
							user_address.address_1,
							user_address.address_2,
							user_address.city_town,
							user_address.province_state,
							user_address.postalcode_zip
						 FROM users 
							LEFT JOIN profile ON users.id = profile.user_id 
							LEFT JOIN user_address ON user_address.profile_id = profile.id
							WHERE $1=$2
						LIMIT 1`

	UpdateStatusQuery = "UPDATE users SET confirmed=TRUE WHERE id=$1"

	createQuery = `INSERT INTO users(first_name, middle_name, last_name, email, country, uuid, created_at) values($1, $2, $3, $4, $5, $6, $7)`

	createContactQuery = `INSERT INTO contacts(first_name, middle_name, last_name, email, mobile, mobile_account, user_id, created_at) values($1, $2, $3, $4, $5, $6, $7, $8)`

	getContactsQuery = "SELECT * FROM contacts WHERE user_id=$1"
)

type PostgService struct {
	DB *sqlx.DB
}

func NewPostgService(db *sqlx.DB) UsersService {
	return &PostgService{DB: db}
}

func (s *PostgService) GetUserByID(ctx context.Context, id int64) (*pb.User, error) {
	u := &pb.User{Id: id}

	if err := s.getUserBy(ctx, "id", id, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *PostgService) GetUserByUUID(ctx context.Context, uuid string) (*pb.User, error) {
	u := &pb.User{Uuid: uuid}

	if err := s.getUserBy(ctx, "uuid", uuid, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *PostgService) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	u := &pb.User{Email: email}

	if err := s.getUserBy(ctx, "email", email, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *PostgService) getUserBy(_ context.Context, kind interface{}, value interface{}, u *pb.User) error {
	//err := s.DB.Get(u, getUserByQuery, kind, value)

	row := s.DB.QueryRow(getUserByQuery, kind, value)

	err := row.Scan(
		&u.Id, &u.Uuid, &u.FirstName, &u.MiddleName, &u.LastName, &u.Email, &u.Confirmed, &u.CreatedAt, &u.Country, &u.Profile.BirthDate, &u.Profile.Gender, &u.Profile.Occupation,
		&u.Profile.Address.Country, &u.Profile.Address.Address1, &u.Profile.Address.Address2, &u.Profile.Address.CityTown, &u.Profile.Address.ProvinceState, &u.Profile.Address.PostalcodeZip,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return err
	}

	return nil
}

func (s *PostgService) UpdateEmail(_ context.Context, u *pb.User) error {
	_, err := s.DB.Exec(UpdateEmailQuery, u.Email, u.Id)

	return err
}

func (s *PostgService) UpdateStatus(_ context.Context, u *pb.User) error {
	_, err := s.DB.Exec(UpdateStatusQuery, u.Id)

	return err
}

func (s *PostgService) Create(_ context.Context, u *pb.User) error {
	_, err := s.DB.Exec(createQuery, u.FirstName, u.FirstName, u.LastName, u.Email, u.Country, u.Uuid, time.Now())

	return err
}

func (s *PostgService) CreateContact(_ context.Context, c *pb.Contact) error {
	_, err := s.DB.Exec(createContactQuery, c.FirstName, c.MiddleName, c.LastName, c.Email, c.Mobile, c.MobileAccount, c.UserId, c.CreatedAt)

	return err
}

func (s *PostgService) GetContacts(_ context.Context, uid string) ([]*pb.Contact, error) {
	var contacts []*pb.Contact

	err := s.DB.Select(contacts, getContactsQuery, uid)

	return contacts, err
}
