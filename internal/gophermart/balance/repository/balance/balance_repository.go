package balance

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/repository/withdrawal"
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/requests"
	"github.com/anoriar/gophermart/internal/gophermart/balance/entity"
	withdrawalRepositoryPkg "github.com/anoriar/gophermart/internal/gophermart/balance/repository/withdrawal"
	"github.com/anoriar/gophermart/internal/gophermart/shared/app/db"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type BalanceRepository struct {
	db                   *db.Database
	withdrawalRepository withdrawalRepositoryPkg.WithdrawalRepositoryInterface
}

func NewBalanceRepository(db *db.Database, withdrawalRepository withdrawalRepositoryPkg.WithdrawalRepositoryInterface) *BalanceRepository {
	return &BalanceRepository{db: db, withdrawalRepository: withdrawalRepository}
}

func (repository BalanceRepository) createBalance(tx *sqlx.Tx, userID string, sum float64) error {
	_, err := tx.NamedExec(
		`INSERT INTO balances (id, user_id, balance, withdrawal, updated_at) 
			VALUES (:id, :user_id, :balance, 0.0, CURRENT_TIMESTAMP(3))`,
		map[string]interface{}{
			"id":      uuid.NewString(),
			"balance": sum,
			"user_id": userID,
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

func (repository BalanceRepository) AddUserBalance(ctx context.Context, userID string, sum float64) error {
	tx, err := repository.db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddUserBalance: %w: %v", errors2.ErrInternalError, err)
	}
	defer tx.Rollback()

	var balance entity.Balance
	err = tx.GetContext(ctx, &balance, "SELECT * FROM balances WHERE user_id = $1 FOR UPDATE", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = repository.createBalance(tx, userID, sum)
			if err != nil {
				return fmt.Errorf("AddUserBalance: %w: %v", errors2.ErrInternalError, err)
			}
		} else {
			return fmt.Errorf("AddUserBalance: %w: %v", errors2.ErrInternalError, err)
		}
	} else {
		_, err = tx.NamedExec(
			`UPDATE balances 
				SET balance = balance + :add_sum, updated_at = CURRENT_TIMESTAMP(3) 
                WHERE user_id = :user_id`,
			map[string]interface{}{
				"add_sum": sum,
				"user_id": userID,
			},
		)
		if err != nil {
			return fmt.Errorf("AddUserBalance: %w: %v", errors2.ErrInternalError, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("AddUserBalance: %w: %v", errors2.ErrInternalError, err)
	}
	return nil
}

func (repository BalanceRepository) WithdrawUserBalance(ctx context.Context, userID string, withdrawDto requests.WithdrawDto) error {
	tx, err := repository.db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("WithdrawUserBalance: %w: %v", errors2.ErrInternalError, err)
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(
		`UPDATE balances 
				SET balance = balance - :withdraw_sum, withdrawal = withdrawal + :withdraw_sum, updated_at = CURRENT_TIMESTAMP(3) 
                WHERE user_id = :user_id`,
		map[string]interface{}{
			"withdraw_sum": withdrawDto.Sum,
			"user_id":      userID,
		},
	)
	if err != nil {
		return fmt.Errorf("WithdrawUserBalance: %w: %v", errors2.ErrInternalError, err)
	}

	err = repository.withdrawalRepository.CreateWithdrawal(ctx, tx, withdrawal.CreateWithdrawalDto{
		UserID: userID,
		Order:  withdrawDto.Order,
		Sum:    withdrawDto.Sum,
	})
	if err != nil {
		return fmt.Errorf("WithdrawUserBalance: %w: %v", errors2.ErrInternalError, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("WithdrawUserBalance: %w: %v", errors2.ErrInternalError, err)
	}
	return nil
}
