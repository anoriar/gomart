package balance

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	balanceDtoPkg "github.com/anoriar/gophermart/internal/gophermart/balance/dto/repository/balance"
	"github.com/anoriar/gophermart/internal/gophermart/balance/entity"
	"github.com/anoriar/gophermart/internal/gophermart/shared/app/db"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type BalanceRepository struct {
	db *db.Database
}

func NewBalanceRepository(db *db.Database) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (repository BalanceRepository) createBalance(tx *sqlx.Tx, updateDto balanceDtoPkg.UpdateBalanceDto) error {
	_, err := tx.NamedExec(
		`INSERT INTO balances (id, user_id, balance, withdrawal, updated_at) 
			VALUES (:id, :user_id, :balance, :withdrawal, CURRENT_TIMESTAMP(3))`,
		map[string]interface{}{
			"id":         uuid.NewString(),
			"balance":    updateDto.Balance,
			"withdrawal": updateDto.Withdrawal,
			"user_id":    updateDto.UserID,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return fmt.Errorf("CreateBalance: %w: %v", errors2.ErrConflict, err)
		}
		return fmt.Errorf("CreateBalance: %w: %v", errors2.ErrInternalError, err)
	}
	return nil
}

func (repository BalanceRepository) updateBalance(tx *sqlx.Tx, updateDto balanceDtoPkg.UpdateBalanceDto) error {
	_, err := tx.NamedExec(
		`UPDATE balances 
				SET balance = :balance, withdrawal = :withdrawal, updated_at = CURRENT_TIMESTAMP(3) 
                WHERE user_id = :user_id`,
		map[string]interface{}{
			"balance":    updateDto.Balance,
			"withdrawal": updateDto.Withdrawal,
			"user_id":    updateDto.UserID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("UpdateBalance: %w: %v", errors2.ErrNotFound, err)
		}
		return fmt.Errorf("UpdateBalance: %w: %v", errors2.ErrInternalError, err)
	}
	return nil
}

func (repository BalanceRepository) UpsertBalance(ctx context.Context, userID string, calcFunc func(curBalance *entity.Balance) balanceDtoPkg.UpdateBalanceDto) error {

	tx, err := repository.db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("UpsertBalance: %w: %v", errors2.ErrInternalError, err)
	}
	defer tx.Rollback()

	var balance entity.Balance
	err = tx.GetContext(ctx, &balance, "SELECT * FROM balances WHERE user_id = $1 FOR UPDATE", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			updateDto := calcFunc(nil)
			err = repository.createBalance(tx, updateDto)
			if err != nil {
				return fmt.Errorf("UpsertBalance: %w: %v", errors2.ErrInternalError, err)
			}
		} else {
			return fmt.Errorf("UpsertBalance: %w: %v", errors2.ErrInternalError, err)
		}
	} else {
		updateDto := calcFunc(&balance)
		err = repository.updateBalance(tx, updateDto)
		if err != nil {
			return fmt.Errorf("UpsertBalance: %w: %v", errors2.ErrInternalError, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UpsertBalance: %w: %v", errors2.ErrInternalError, err)
	}
	return nil
}

func (repository BalanceRepository) GetBalanceByUserID(ctx context.Context, userID string) (entity.Balance, error) {
	var resultBalance entity.Balance
	err := repository.db.Conn.QueryRowxContext(ctx, "SELECT * FROM balances WHERE user_id = $1", userID).StructScan(&resultBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Balance{}, fmt.Errorf("GetBalanceByUserID: %w: %v", errors2.ErrNotFound, err)
		}
		return entity.Balance{}, fmt.Errorf("GetBalanceByUserID: %w: %v", errors2.ErrInternalError, err)
	}
	return resultBalance, nil
}
