package balance

import (
	"context"
	balanceResponsePkg "github.com/anoriar/gophermart/internal/gophermart/dto/responses/balance"
)

type BalanceServiceInterface interface {
	GetUserBalance(ctx context.Context, userID string) (balanceResponsePkg.BalanceResponseDto, error)
	SyncUserBalance(ctx context.Context, userID string) error
}
