package withdraw

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/requests"
	"github.com/anoriar/gophermart/internal/gophermart/balance/entity"
)

type WithdrawServiceInterface interface {
	Withdraw(ctx context.Context, userID string, withdraw requests.WithdrawDto) error
	GetWithdrawalsByUserID(ctx context.Context, userID string) ([]entity.Withdrawal, error)
}
