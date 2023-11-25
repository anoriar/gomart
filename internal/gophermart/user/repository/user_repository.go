package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/shared/app/db"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/shared/errors"
	"github.com/anoriar/gophermart/internal/gophermart/user/entity"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	db *db.Database
}

func NewUserRepository(db *db.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository *UserRepository) AddUser(ctx context.Context, user entity.User) error {
	_, err := repository.db.Conn.ExecContext(ctx, "INSERT INTO users (id, login, password, salt) VALUES ($1, $2, $3, $4)", user.ID, user.Login, user.Password, user.Salt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return fmt.Errorf("%w: %v", errors2.ErrConflict, err)
		}
		return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}
	return nil
}

func (repository *UserRepository) GetUserByLogin(ctx context.Context, login string) (entity.User, error) {

	var userRes entity.User
	err := repository.db.Conn.QueryRowxContext(ctx, "SELECT id, login, password, salt FROM users WHERE login=$1", login).StructScan(&userRes)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userRes, fmt.Errorf("%w: %v", errors2.ErrNotFound, err)
		}
		return userRes, fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}
	return userRes, nil
}
