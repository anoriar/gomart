package withdrawal

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/entity/withdrawal"
)

type WithdrawalRepositoryInterface interface {
	CreateWithdrawal(ctx context.Context, withdrawal withdrawal.Withdrawal) error
	GetWithdrawalsByUserID(ctx context.Context, userID string) ([]withdrawal.Withdrawal, error)
}
