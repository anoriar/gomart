package db

import "github.com/jmoiron/sqlx"

type Database struct {
	Conn *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{Conn: db}
}

func (db *Database) Ping() error {
	err := db.Conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Close() {
	db.Conn.Close()
}
