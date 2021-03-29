package service

import (
	"context"
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

// UsersService describes the User service.
type UsersService interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByUUID(ctx context.Context, uuid string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateEmail(ctx context.Context, user *User) error
	CreateContact(ctx context.Context, contact *Contact) (*Contact, error)
	UpdateContact(ctx context.Context, contact *Contact) (*Contact, error)
	GetContacts(ctx context.Context, uid int64) ([]*Contact, error)
	DeleteContact(ctx context.Context, contact *Contact) error
	SetUserProfile(ctx context.Context, user *User) (*User, error)
	UpdateUserProfile(ctx context.Context, user *User) (*User, error)
}

type User struct {
	FirstName  string     `pjson:"firstName,omitempty" db:"first_name"`
	MiddleName string     `json:"middleName,omitempty" db:"middle_name"`
	LastName   string     `json:"lastName,omitempty" db:"last_name"`
	Email      string     `json:"email,omitempty" db:"email"`
	Uuid       string     `json:"uuid,omitempty" db:"uuid"`
	Id         int64      `json:"id,omitempty" db:"id"`
	Country    string     `json:"country,omitempty" db:"country"`
	CreatedAt  *time.Time `json:"createdAt,omitempty" db:"created_at"`
	Confirmed  bool       `json:"confirmed,omitempty" db:"confirmed"`
	Profile    *Profile   `json:"profile,omitempty"`
}

type Profile struct {
	Gender     string     `json:"gender,omitempty" db:"gender"`
	Occupation string     `json:"occupation,omitempty" db:"occupation"`
	BirthDate  *time.Time `json:"birthDate,omitempty" db:"birth_date"`
	Address    *Address   `json:"address,omitempty"`
}

type Address struct {
	Country       string `json:"country,omitempty" db:"country"`
	Address1      string `json:"address1,omitempty" db:"address_1"`
	Address2      string `json:"address2,omitempty" db:"address_2"`
	CityTown      string `json:"cityTown,omitempty" db:"city_town"`
	ProvinceState string `json:"provinceState,omitempty" db:"province_state"`
	PostalcodeZip string `json:"postalcodeZip,omitempty" db:"postal_code_zip"`
}

type Contact struct {
	FirstName     string     `json:"firstName,omitempty"`
	MiddleName    string     `json:"middleName,omitempty"`
	LastName      string     `json:"lastName,omitempty"`
	Email         string     `json:"email,omitempty"`
	Mobile        string     `json:"mobile,omitempty"`
	MobileAccount string     `json:"mobileAccount,omitempty"`
	Id            int64      `json:"id,omitempty"`
	UserId        string     `json:"userId,omitempty"`
	CreatedAt     *time.Time `json:"createdAt,omitempty"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

// New returns a UsersService with all of the expected middleware wired in.
func New(svc UsersService, middleware []Middleware) UsersService {
	for _, m := range middleware {
		svc = m(svc)
	}

	return svc
}
