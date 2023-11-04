package user

import (
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"go.uber.org/zap"
)

type UserRepository struct {
	db     *db.Database
	logger *zap.Logger
}

func NewUserRepository(db *db.Database, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}
