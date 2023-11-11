package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/entity/user"
	"github.com/anoriar/gophermart/internal/gophermart/repository/repository_error"
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

func (repository *UserRepository) AddUser(ctx context.Context, user user.User) error {
	_, err := repository.db.Conn.ExecContext(ctx, "INSERT INTO users (id, login, password, salt) VALUES ($1, $2, $3, $4)", user.Id, user.Login, user.Password, user.Salt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return repository_error.ErrConflict
		}
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (repository *UserRepository) GetUserByLogin(ctx context.Context, login string) (user.User, error) {

	var userRes user.User
	err := repository.db.Conn.QueryRowxContext(ctx, "SELECT id, login, password, salt FROM users WHERE login=$1", login).StructScan(&userRes)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userRes, repository_error.ErrNotFound
		}
		return userRes, fmt.Errorf("%w", err)
	}
	return userRes, nil
}
