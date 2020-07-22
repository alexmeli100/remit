package service

import (
	"context"
	"fmt"
	"github.com/alexmeli100/remit/users/pkg/grpc/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var pg PostgService

const TableCreationQuery = `CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL,
    uuid       TEXT    NOT NULL,
    first_name  TEXT    NOT NULL,
    last_name   TEXT    NOT NULL,
    email      TEXT    NOT NULL,
    confirmed  BOOLEAN NOT NULL,
    country    Text    Not NuLL,
    CONSTRAINT user_pkey PRIMARY KEY (id)
)`

func openConnection() (*sqlx.DB, error) {
	pass := os.Getenv("POSTGRES_PASSWORD")
	userName := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	connString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", userName, pass, dbName)

	db, err := sqlx.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func ensureTableExists() {
	if _, err := pg.DB.Exec(TableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	if _, err := pg.DB.Exec("DELETE FROM users"); err != nil {
		log.Fatal(err)
	}

	if _, err := pg.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1"); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	db, err := openConnection()

	if err != nil {
		log.Fatal(err)
	}

	pg = PostgService{db}
	ensureTableExists()
	code := m.Run()
	clearTable()
	pg.DB.Close()
	os.Exit(code)

}

func compare(u1 *pb.User, u2 *pb.User) bool {
	return proto.Equal(u1, u2)
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

	if err := pg.Create(context.Background(), u); err != nil {
		t.Errorf("%v", err)
	}

	tu, err := pg.GetUserByUUID(context.Background(), uid)

	if err != nil {
		t.Errorf("%v", err)
	}

	if !compare(u, tu) {
		t.Errorf("Expected: %v, Got: %v", u, tu)
	}
}

func TestPostgService_GetUserByID(t *testing.T) {
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

	if err := pg.Create(context.Background(), u); err != nil {
		t.Errorf("%v", err)
	}

	tu, err := pg.GetUserByUUID(context.Background(), uid)

	if err != nil {
		t.Errorf("%v", err)
	}

	if !compare(u, tu) {
		t.Errorf("Expected: %v, Got: %v", u, tu)
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

	if err := pg.Create(context.Background(), u); err != nil {
		t.Errorf("%v", err)
	}

	tu, err := pg.GetUserByEmail(context.Background(), u.Email)

	if err != nil {
		t.Errorf("%v", err)
	}

	if !compare(u, tu) {
		t.Errorf("Expected: %v, Got: %v", u, tu)
	}
}

func TestLoggingMiddleware_GetUserByID(t *testing.T) {
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

	if err := pg.Create(context.Background(), u); err != nil {
		t.Errorf("%v", err)
	}

	tu, err := pg.GetUserByID(context.Background(), 1)

	if err != nil {
		t.Errorf("%v", err)
	}

	if !compare(u, tu) {
		t.Errorf("Expected: %v, Got: %v", u, tu)
	}
}
