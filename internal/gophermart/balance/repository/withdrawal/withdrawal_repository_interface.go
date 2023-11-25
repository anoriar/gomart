package withdrawal

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/balance/entity"
)

type WithdrawalRepositoryInterface interface {
	CreateWithdrawal(ctx context.Context, withdrawal entity.Withdrawal) error
	GetWithdrawalsByUserID(ctx context.Context, userID string) ([]entity.Withdrawal, error)
}
