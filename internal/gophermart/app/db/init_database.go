package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"log"
)

func InitializeDatabase(dsn string) (*Database, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = goose.Up(db.DB, "migrations")
	if err != nil {
		log.Fatal(err)
	}
	return NewDatabase(db), nil
}
