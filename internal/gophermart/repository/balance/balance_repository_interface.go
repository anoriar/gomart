package balance

import (
	"context"
	balanceDtoPkg "github.com/anoriar/gophermart/internal/gophermart/dto/repository/balance"
	"github.com/anoriar/gophermart/internal/gophermart/entity/balance"
)

type BalanceRepositoryInterface interface {
	UpsertBalance(ctx context.Context, userID string, calcFunc func(curBalance *balance.Balance) balanceDtoPkg.UpdateBalanceDto) error
	GetBalanceByUserID(ctx context.Context, userID string) (balance.Balance, error)
}
