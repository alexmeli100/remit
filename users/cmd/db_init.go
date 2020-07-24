package main

import "github.com/jmoiron/sqlx"

var tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
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

func initDB(db *sqlx.DB) error {
	_, err := db.Exec(tableCreationQuery)
	return err
}
