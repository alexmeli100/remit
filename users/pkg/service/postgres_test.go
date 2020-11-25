package service

import (
	"context"
	"fmt"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

var pg PostgService

func openConnection() (*sqlx.DB, error) {
	pass := "ogttpy3ikvctrlhs"
	userName := "doadmin"
	host := "db-postgresql-nyc1-31181-do-user-8015056-0.b.db.ondigitalocean.com"
	port := "25060"
	dbName := "wealow-users-test"
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, userName, pass, dbName)

	db, err := sqlx.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func clearTable() {
	if _, err := pg.DB.Exec("DELETE FROM users"); err != nil {
		log.Fatal(err)
	}
}

func TestPostGres(t *testing.T) {
	db, err := openConnection()

	if err != nil {
		log.Fatal(err)
	}

	pg = PostgService{db}

	t.Run("TestPostGres", func(t *testing.T) {
		t.Run("TestCreate", func(t *testing.T) {
			TestPostgService_Create(t)
		})

		t.Run("TestGetUserByUUID", func(t *testing.T) {
			TestPostgService_GetUserByUUID(t)
		})

		t.Run("TestGetUserByEmail", func(t *testing.T) {
			TestPostgService_GetUserByEmail(t)
		})

		t.Run("TestUpdateEmail", func(t *testing.T) {
			TestPostgService_UpdateEmail(t)
		})

		t.Run("TestSetUSerProfile", func(t *testing.T) {
			TestPostgService_SetUserProfile(t)
		})

		t.Run("TestUpdateUserProfile", func(t *testing.T) {
			TestPostgService_UpdateUserProfile(t)
		})
	})
}

func TestPostgService_Create(t *testing.T) {
	clearTable()

	uid := uuid.New().String()
	u := &pb.User{
		FirstName: "Alex",
		LastName:  "Meli",
		Uuid:      uid,
		Email:     "alexmeli100@gmail.com",
		Confirmed: false,
		Country:   "Canada",
	}

	user, err := pg.Create(context.Background(), u)

	if assert.NoError(t, err) {
		assert.NotNil(t, user)
	}
}

func TestPostgService_GetUserByUUID(t *testing.T) {
	clearTable()

	uid := uuid.New().String()
	u := &pb.User{
		FirstName: "James",
		LastName:  "Meli",
		Uuid:      uid,
		Email:     "jamesmeli100@gmail.com",
		Confirmed: false,
		Country:   "USA",
	}

	user, err := pg.Create(context.Background(), u)

	if assert.NoError(t, err) {
		assert.NotNil(t, user)
	}

	tu, err := pg.GetUserByUUID(context.Background(), uid)

	if assert.NoError(t, err) {
		assert.NotNil(t, tu)
	}
}

func TestPostgService_GetUserByEmail(t *testing.T) {
	clearTable()

	uid := uuid.New().String()
	u := &pb.User{
		FirstName: "James",
		LastName:  "Meli",
		Uuid:      uid,
		Email:     "jamesmeli100@gmail.com",
		Confirmed: false,
		Country:   "USA",
	}

	user, err := pg.Create(context.Background(), u)

	if assert.NoError(t, err) {
		assert.NotNil(t, user)
	}

	tu, err := pg.GetUserByEmail(context.Background(), user.Email)

	if assert.NoError(t, err) {
		assert.NotNil(t, tu)
	}
}

func TestPostgService_UpdateEmail(t *testing.T) {
	clearTable()

	uid := uuid.New().String()
	u := &pb.User{
		FirstName: "James",
		LastName:  "Meli",
		Uuid:      uid,
		Email:     "jamesmeli100@gmail.com",
		Confirmed: false,
		Country:   "USA",
	}

	user, err := pg.Create(context.Background(), u)

	if assert.NoError(t, err) {
		assert.NotNil(t, user)
	}

	newUser := &pb.User{
		Id:    user.Id,
		Email: "alexmeli100@gmail.com",
	}

	err = pg.UpdateEmail(context.Background(), newUser)

	assert.NoError(t, err)
}

func TestPostgService_SetUserProfile(t *testing.T) {
	clearTable()

	uid := uuid.New().String()
	u := &pb.User{
		FirstName: "James",
		LastName:  "Meli",
		Uuid:      uid,
		Email:     "jamesmeli100@gmail.com",
		Confirmed: false,
		Country:   "USA",
	}

	user, err := pg.Create(context.Background(), u)

	if assert.NoError(t, err) {
		assert.NotNil(t, user)
	}

	b := time.Now()

	p := &pb.Profile{
		BirthDate:  &b,
		Occupation: "Student",
		Gender:     "Male",
		Address: &pb.Address{
			Address1:      "110 York Mills Rd",
			Address2:      "",
			Country:       "Canada",
			CityTown:      "Toronto",
			ProvinceState: "On",
			PostalcodeZip: "M2N13J",
		},
	}

	user.Profile = p
	//t.Log(user)

	updatedUser, err := pg.SetUserProfile(context.Background(), user)

	if assert.NoError(t, err) {
		assert.NotNil(t, updatedUser)
	}
}

func TestPostgService_UpdateUserProfile(t *testing.T) {
	clearTable()

	uid := uuid.New().String()
	u := &pb.User{
		FirstName: "James",
		LastName:  "Meli",
		Uuid:      uid,
		Email:     "jamesmeli100@gmail.com",
		Confirmed: false,
		Country:   "USA",
	}

	user, err := pg.Create(context.Background(), u)

	if assert.NoError(t, err) {
		assert.NotNil(t, user)
	}

	b := time.Now()

	p := &pb.Profile{
		BirthDate:  &b,
		Occupation: "Student",
		Gender:     "Male",
		Address: &pb.Address{
			Address1:      "110 York Mills Rd",
			Address2:      "",
			Country:       "Canada",
			CityTown:      "Toronto",
			ProvinceState: "On",
			PostalcodeZip: "M2N13J",
		},
	}

	user.Profile = p

	updatedUser, err := pg.SetUserProfile(context.Background(), user)

	if assert.NoError(t, err) {
		assert.NotNil(t, updatedUser)
	}

	newP := &pb.Profile{
		BirthDate:  &b,
		Occupation: "Doctor",
		Gender:     "Female",
		Address: &pb.Address{
			Address1:      "2900 Jane Street",
			Address2:      "",
			Country:       "Canada",
			CityTown:      "Toronto",
			ProvinceState: "On",
			PostalcodeZip: "M2N13J",
		},
	}

	user.Profile = newP

	updatedUser, err = pg.UpdateUserProfile(context.Background(), user)

	if assert.NoError(t, err) {
		assert.NotNil(t, updatedUser)
	}
}
