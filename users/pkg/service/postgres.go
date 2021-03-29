package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const (
	UpdateEmailQuery = "UPDATE users SET email=$1 WHERE id=$2"

	getUserByQuery = `
		SELECT
			users.id, users.uuid,users.first_name, users.middle_name, users.last_name, users.email,users.confirmed, 
            users.created_at, users.country, profile.birth_date, COALESCE(profile.gender, ''),COALESCE(profile.occupation, ''),
            COALESCE(user_address.country, ''), COALESCE(user_address.address_1, ''), COALESCE(user_address.address_2, ''),
            COALESCE(user_address.city_town, ''), COALESCE(user_address.province_state, ''), COALESCE(user_address.postalcode_zip, '')
		FROM users
		LEFT JOIN profile ON users.id = profile.user_id
		LEFT JOIN user_address ON user_address.profile_id = profile.id
		WHERE %s=$1
		LIMIT 1`

	createQuery = `
		INSERT INTO users(first_name, middle_name, last_name, email, country, uuid, created_at) 
		values($1, $2, $3, $4, $5, $6, $7) RETURNING *`

	createContactQuery = `
		INSERT INTO contacts(first_name, middle_name, last_name, email, mobile, mobile_account, user_id, created_at) 
		values($1, $2, $3, $4, $5, $6, $7, $8)`

	updateContactQuery = `
		UPDATE contacts SET first_name = $1, last_name = $2, email = $3, mobile = $4, mobile_account = $5 
		WHERE id = $6
		RETURNING *; 
	`

	createUserProfileQuery = `
		with updated_profile as (
    		INSERT INTO profile(birth_date, user_id, gender, occupation)
    		VALUES($1, $2, $3, $4) RETURNING id, birth_date, gender, occupation)
		INSERT INTO user_address(profile_id, country, address_1, address_2, city_town, province_state, postalcode_zip)
    		VALUES ((select id FROM updated_profile), $5, $6, $7, $8, $9, $10)
    		RETURNING 
    			(select birth_date FROM updated_profile), (select gender FROM updated_profile), (select occupation FROM updated_profile),
    			country, address_1, address_2, city_town, province_state, postalcode_zip;
	`

	updateUserProfileQuery = `
		with updated_profile as (
    		UPDATE profile SET birth_date = $1, gender = $2, occupation = $3
    		WHERE user_id = $4 RETURNING birth_date, gender, occupation)
		UPDATE user_address SET country = $5, address_1 = $6, address_2 = $7, city_town = $8, province_state = $9, postalcode_zip = $10
		WHERE profile_id = (SELECT id from updated_profile) 
		RETURNING
			(select birth_date FROM updated_profile), (select gender FROM updated_profile), (select occupation FROM updated_profile),
			country, address_1, address_2, city_town, province_state, postalcode_zip;
	`

	deleteContactQuery = `DELETE FROM contacts where id=$1`

	getContactsQuery = "SELECT * FROM contacts WHERE user_id=$1"
)

type PostgService struct {
	DB *sqlx.DB
}

func NewPostgService(db *sqlx.DB) UsersService {
	return &PostgService{DB: db}
}

func (s *PostgService) GetUserByID(ctx context.Context, id int64) (*User, error) {
	u := &User{
		Id:      id,
		Profile: &Profile{Address: &Address{}},
	}

	if err := s.getUserBy(ctx, "id", id, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *PostgService) GetUserByUUID(ctx context.Context, uuid string) (*User, error) {
	u := &User{
		Uuid:    uuid,
		Profile: &Profile{Address: &Address{}},
	}

	if err := s.getUserBy(ctx, "uuid", uuid, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *PostgService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{
		Email:   email,
		Profile: &Profile{Address: &Address{}},
	}

	if err := s.getUserBy(ctx, "email", email, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *PostgService) getUserBy(ctx context.Context, kind string, value interface{}, u *User) error {
	//err := s.DB.Get(u, getUserByQuery, kind, value)
	q := fmt.Sprintf(getUserByQuery, "users."+kind)
	row := s.DB.QueryRowContext(ctx, q, value)

	err := row.Scan(
		&u.Id, &u.Uuid, &u.FirstName, &u.MiddleName, &u.LastName, &u.Email, &u.Confirmed,
		&u.CreatedAt, &u.Country, &u.Profile.BirthDate, &u.Profile.Gender, &u.Profile.Occupation,
		&u.Profile.Address.Country, &u.Profile.Address.Address1, &u.Profile.Address.Address2,
		&u.Profile.Address.CityTown, &u.Profile.Address.ProvinceState, &u.Profile.Address.PostalcodeZip,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		return err
	}

	return nil
}

func scanProfile(row *sqlx.Row, p *Profile) error {
	err := row.Scan(
		&p.BirthDate, &p.Gender, &p.Occupation, &p.Address.Country, &p.Address.Address1,
		&p.Address.Address2, &p.Address.CityTown, &p.Address.ProvinceState, &p.Address.PostalcodeZip)

	return err
}

func (s *PostgService) UpdateEmail(ctx context.Context, u *User) error {
	_, err := s.DB.ExecContext(ctx, UpdateEmailQuery, u.Email, u.Id)

	return err
}

func (s *PostgService) CreateUser(ctx context.Context, u *User) (*User, error) {
	user := &User{}

	err := s.DB.QueryRowxContext(
		ctx,
		createQuery,
		u.FirstName, u.MiddleName, u.LastName, u.Email, u.Country, u.Uuid, time.Now(),
	).StructScan(user)

	return user, err
}

func (s *PostgService) SetUserProfile(ctx context.Context, u *User) (*User, error) {
	row := s.DB.QueryRowxContext(
		ctx,
		createUserProfileQuery,
		u.Profile.BirthDate, u.Id, u.Profile.Gender, u.Profile.Occupation,
		u.Profile.Address.Country, u.Profile.Address.Address1, u.Profile.Address.Address2,
		u.Profile.Address.CityTown, u.Profile.Address.ProvinceState, u.Profile.Address.PostalcodeZip)

	p := &Profile{Address: &Address{}}

	if err := scanProfile(row, p); err != nil {
		return nil, err
	}

	u.Profile = p
	return u, nil
}

func (s *PostgService) UpdateUserProfile(ctx context.Context, u *User) (*User, error) {
	row := s.DB.QueryRowxContext(
		ctx,
		updateUserProfileQuery,
		u.Profile.BirthDate, u.Profile.Gender, u.Profile.Occupation, u.Id,
		u.Profile.Address.Country, u.Profile.Address.Address1, u.Profile.Address.Address2,
		u.Profile.Address.CityTown, u.Profile.Address.ProvinceState, u.Profile.Address.PostalcodeZip)

	p := &Profile{Address: &Address{}}

	if err := scanProfile(row, p); err != nil {
		return nil, err
	}

	u.Profile = p
	return u, nil
}

func (s *PostgService) CreateContact(ctx context.Context, c *Contact) (*Contact, error) {
	contact := &Contact{}

	err := s.DB.QueryRowxContext(
		ctx,
		createContactQuery,
		c.FirstName, c.MiddleName, c.LastName, c.Email, c.Mobile, c.MobileAccount, c.UserId, c.CreatedAt,
	).StructScan(contact)

	if err != nil {
		return nil, err
	}

	return contact, err
}

func (s *PostgService) UpdateContact(ctx context.Context, c *Contact) (*Contact, error) {
	contact := &Contact{}
	err := s.DB.QueryRowxContext(
		ctx,
		updateContactQuery,
		c.FirstName, c.LastName, c.Email, c.Mobile, c.MobileAccount, c.Id,
	).StructScan(contact)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *PostgService) DeleteContact(ctx context.Context, contact *Contact) error {
	_, err := s.DB.ExecContext(ctx, deleteContactQuery, contact.Id)

	return err
}

func (s *PostgService) GetContacts(ctx context.Context, uid int64) ([]*Contact, error) {
	var contacts []*Contact

	err := s.DB.SelectContext(ctx, contacts, getContactsQuery, uid)

	return contacts, err
}
