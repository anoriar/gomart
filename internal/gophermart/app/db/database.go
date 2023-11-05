package db

import "github.com/jmoiron/sqlx"

type Database struct {
	conn *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{conn: db}
}

func (db *Database) Ping() error {
	err := db.conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Close() {
	db.Close()
}
