package db

//go:generate mockgen -source=database_interface.go -destination=mock/database.go -package=mock
type DatabaseInterface interface {
	Ping() error
	Close()
}
