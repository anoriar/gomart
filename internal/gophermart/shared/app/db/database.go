package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type Database struct {
	Conn *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{Conn: db}
}

func (db *Database) Ping(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Database::Ping")
	defer span.Finish()

	err := db.Conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Close() {
	db.Conn.Close()
}
