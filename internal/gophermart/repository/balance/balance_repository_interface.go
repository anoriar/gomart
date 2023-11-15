package balance

import (
	"context"
	balanceDtoPkg "github.com/anoriar/gophermart/internal/gophermart/dto/repository/balance"
	"github.com/anoriar/gophermart/internal/gophermart/entity/balance"
)

type BalanceRepositoryInterface interface {
	CreateBalance(ctx context.Context, updateDto balanceDtoPkg.UpdateBalanceDto) error
	UpdateBalance(ctx context.Context, updateDto balanceDtoPkg.UpdateBalanceDto) error
	GetBalanceByUserID(ctx context.Context, userID string) (balance.Balance, error)
}
