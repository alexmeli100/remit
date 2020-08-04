package main

import (
	"fmt"
	"github.com/alexmeli100/remit/users/cmd/service"
	userService "github.com/alexmeli100/remit/users/pkg/service"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

func main() {
	password := os.Getenv("POSTGRES_PASSWORD")
	userName := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("USER_DB_SERVICE_HOST")
	dbPort := os.Getenv("USER_DB_SERVICE_PORT")
	db, err := openDB(dbHost, dbPort, userName, password, dbName)

	if err != nil {
		log.Fatal(err)
	}

	if err = initDB(db); err != nil {
		log.Fatal(err)
	}

	pg := userService.NewPostgService(db)
	service.Run(pg)
}

func openDB(host, port, userName, password, dbName string) (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, userName, password, dbName)

	db, err := sqlx.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
