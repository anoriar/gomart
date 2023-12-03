package withdrawal

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/repository/withdrawal"
	"github.com/anoriar/gophermart/internal/gophermart/balance/entity"
	"github.com/jmoiron/sqlx"
)

type WithdrawalRepositoryInterface interface {
	CreateWithdrawal(ctx context.Context, tx *sqlx.Tx, createDto withdrawal.CreateWithdrawalDto) error
	GetWithdrawalsByUserID(ctx context.Context, userID string) ([]entity.Withdrawal, error)
}
