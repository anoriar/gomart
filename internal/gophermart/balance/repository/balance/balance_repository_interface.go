package balance

import (
	"context"
	balanceDtoPkg "github.com/anoriar/gophermart/internal/gophermart/balance/dto/repository/balance"
	"github.com/anoriar/gophermart/internal/gophermart/balance/entity"
)

type BalanceRepositoryInterface interface {
	UpsertBalance(ctx context.Context, userID string, calcFunc func(curBalance *entity.Balance) balanceDtoPkg.UpdateBalanceDto) error
	GetBalanceByUserID(ctx context.Context, userID string) (entity.Balance, error)
}
