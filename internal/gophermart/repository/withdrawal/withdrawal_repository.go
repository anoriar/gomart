package withdrawal

import (
	"context"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	"github.com/anoriar/gophermart/internal/gophermart/entity/withdrawal"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type WithdrawalRepository struct {
	db *db.Database
}

func NewWithdrawalRepository(db *db.Database) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

func (repository WithdrawalRepository) CreateWithdrawal(ctx context.Context, withdrawal withdrawal.Withdrawal) error {
	_, err := repository.db.Conn.NamedExecContext(ctx,
		`INSERT INTO withdrawals (id, user_id, order_id, sum, processed_at) 
			VALUES (:id, :user_id, :order_id, :sum, CURRENT_TIMESTAMP(3))`,
		withdrawal,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return fmt.Errorf("CreateWithdrawal: %w: %v", domain_errors.ErrConflict, err)
		}
		return fmt.Errorf("CreateWithdrawal: %w: %v", domain_errors.ErrInternalError, err)
	}
	return nil
}

func (repository WithdrawalRepository) GetWithdrawalsByUserID(ctx context.Context, userID string) ([]withdrawal.Withdrawal, error) {
	var resultWithdrawals []withdrawal.Withdrawal
	rows, err := repository.db.Conn.QueryxContext(ctx, "SELECT * FROM withdrawals WHERE user_id = $1 ORDER BY processed_at", userID)
	defer rows.Close()

	for rows.Next() {
		var withdrawal withdrawal.Withdrawal
		err := rows.StructScan(&withdrawal)
		if err != nil {
			return resultWithdrawals, fmt.Errorf("GetWithdrawalsByUserID: %w: %v", domain_errors.ErrInternalError, err)
		}
		resultWithdrawals = append(resultWithdrawals, withdrawal)
	}

	if rows.Err() != nil {
		return resultWithdrawals, fmt.Errorf("GetWithdrawalsByUserID: %w: %v", domain_errors.ErrInternalError, err)
	}

	return resultWithdrawals, nil
}
