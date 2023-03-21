package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

type DB struct {
	Conn *sql.DB
}

func New(connectionString string) *DB {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	pgDriver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		panic(fmt.Errorf("unable to create postgres driver instance: %v", err))
	}

	if err := RunMigrations(pgDriver); err != nil {
		panic(fmt.Errorf("failed to run migrations: %v", err))
	}

	return &DB{Conn: conn}
}
