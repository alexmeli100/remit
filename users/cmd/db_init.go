package main

import "github.com/jmoiron/sqlx"

var tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
	(
	    id          SERIAL,
	    uuid        TEXT 	  UNIQUE NOT NULL,
	    first_name  TEXT      NOT NULL,
	    last_name   TEXT      NOT NULL,
	    email       TEXT 	  UNIQUE NOT NULL,
	    confirmed   BOOLEAN   NOT NULL,
	    created_at  TIMESTAMP NOT NULL,
	    country     Text      Not NuLL,
	    CONSTRAINT  user_pkey PRIMARY KEY (id)
	)`

func initDB(db *sqlx.DB) error {
	_, err := db.Exec(tableCreationQuery)
	return err
}
