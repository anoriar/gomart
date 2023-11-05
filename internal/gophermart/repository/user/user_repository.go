package user

import (
	"context"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/entity/user"
	"github.com/anoriar/gophermart/internal/gophermart/repository/repository_error"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (repository *UserRepository) AddUser(ctx context.Context, user user.User) error {
	_, err := repository.db.Conn.ExecContext(ctx, "INSERT INTO users (id, login, password, salt) VALUES ($1, $2, $3, $4)", user.Id, user.Login, user.Password, user.Salt)
	if err != nil {
		repository.logger.Error(err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return repository_error.ErrConflict
		}
		return err
	}
	return nil
}
