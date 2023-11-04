package db

import "github.com/jmoiron/sqlx"

type Database struct {
	conn *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{conn: db}
}

func (db *Database) Close() {
	db.Close()
}
