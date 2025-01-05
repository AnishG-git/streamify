package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	username      = "postgres"
	password      = "postgres"
	pgcontainer   = "db"
	migrationsDir = "./migrations"
)

func loadDB(dbname string, external bool) (*sql.DB, error) {
	// Load the database
	host := "db"
	port := 5432
	if external {
		host = "localhost"
		port = 5434
	}
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", username, password, dbname, host, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
