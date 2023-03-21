package database

import (
	"database/sql"
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
	return &DB{Conn: conn}
}
