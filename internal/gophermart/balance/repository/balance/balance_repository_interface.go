package balance

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/requests"
	"github.com/anoriar/gophermart/internal/gophermart/balance/entity"
)

type BalanceRepositoryInterface interface {
	GetBalanceByUserID(ctx context.Context, userID string) (entity.Balance, error)
	AddUserBalance(ctx context.Context, userID string, sum float64) error
	WithdrawUserBalance(ctx context.Context, userID string, withdrawDto requests.WithdrawDto) error
}
