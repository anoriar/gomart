package withdrawal

import (
	"context"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/repository/withdrawal"
	"github.com/anoriar/gophermart/internal/gophermart/balance/entity"
	"github.com/anoriar/gophermart/internal/gophermart/shared/app/db"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type WithdrawalRepository struct {
	db *db.Database
}

func NewWithdrawalRepository(db *db.Database) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

func (repository WithdrawalRepository) CreateWithdrawal(ctx context.Context, tx *sqlx.Tx, createDto withdrawal.CreateWithdrawalDto) error {
	_, err := tx.NamedExecContext(ctx,
		`INSERT INTO withdrawals (id, user_id, order_id, sum, processed_at) 
			VALUES (:id, :user_id, :order_id, :sum, CURRENT_TIMESTAMP(3))`,
		map[string]interface{}{
			"id":       uuid.NewString(),
			"user_id":  createDto.UserID,
			"order_id": createDto.Order,
			"sum":      createDto.Sum,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return fmt.Errorf("CreateWithdrawal: %w: %v", errors2.ErrConflict, err)
		}
		return fmt.Errorf("CreateWithdrawal: %w: %v", errors2.ErrInternalError, err)
	}
	return nil
}

func (repository WithdrawalRepository) GetWithdrawalsByUserID(ctx context.Context, userID string) ([]entity.Withdrawal, error) {
	var resultWithdrawals []entity.Withdrawal
	rows, err := repository.db.Conn.QueryxContext(ctx, "SELECT * FROM withdrawals WHERE user_id = $1 ORDER BY processed_at", userID)
	if err != nil {
		return resultWithdrawals, fmt.Errorf("GetWithdrawalsByUserID: %w: %v", errors2.ErrInternalError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var withdrawal entity.Withdrawal
		err := rows.StructScan(&withdrawal)
		if err != nil {
			return resultWithdrawals, fmt.Errorf("GetWithdrawalsByUserID: %w: %v", errors2.ErrInternalError, err)
		}
		resultWithdrawals = append(resultWithdrawals, withdrawal)
	}

	if rows.Err() != nil {
		return resultWithdrawals, fmt.Errorf("GetWithdrawalsByUserID: %w: %v", errors2.ErrInternalError, err)
	}

	return resultWithdrawals, nil
}
