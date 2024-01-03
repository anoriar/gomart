package db

import "golang.org/x/net/context"

//go:generate mockgen -source=database_interface.go -destination=mock/database.go -package=mock
type DatabaseInterface interface {
	Ping(ctx context.Context) error
	Close()
}
