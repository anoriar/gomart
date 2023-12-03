package balance

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/requests"
	balanceResponsePkg "github.com/anoriar/gophermart/internal/gophermart/balance/dto/responses"
)

type BalanceServiceInterface interface {
	GetUserBalance(ctx context.Context, userID string) (balanceResponsePkg.BalanceResponseDto, error)
	AddUserBalance(ctx context.Context, userID string, sum float64) error
	WithdrawUserBalance(ctx context.Context, userID string, withdrawDto requests.WithdrawDto) error
}
