package balance

import (
	"context"
	balanceResponsePkg "github.com/anoriar/gophermart/internal/gophermart/balance/dto/responses"
)

type BalanceServiceInterface interface {
	GetUserBalance(ctx context.Context, userID string) (balanceResponsePkg.BalanceResponseDto, error)
	UpdateUserBalance(ctx context.Context, userID string, addSum float64, withdrawSum float64) error
}
