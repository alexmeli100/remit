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
	db, err := openDB(userName, password, dbName)

	if err != nil {
		log.Fatal(err)
	}

	pg := userService.NewPostgService(db)
	service.Run(pg)
}

func openDB(userName, password, dbName string) (*sqlx.DB, error) {

	connString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", userName, password, dbName)

	db, err := sqlx.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
