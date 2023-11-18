package withdraw

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/withdraw"
	"github.com/anoriar/gophermart/internal/gophermart/entity/withdrawal"
)

type WithdrawServiceInterface interface {
	Withdraw(ctx context.Context, userID string, withdraw withdraw.WithdrawDto) error
	GetWithdrawalsByUserID(ctx context.Context, userID string) ([]withdrawal.Withdrawal, error)
}
