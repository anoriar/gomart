package balance

import (
	"context"
	"errors"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/domainerrors"
	balanceDtoPkg "github.com/anoriar/gophermart/internal/gophermart/dto/repository/balance"
	balanceResponsePkg "github.com/anoriar/gophermart/internal/gophermart/dto/responses/balance"
	balanceEntityPkg "github.com/anoriar/gophermart/internal/gophermart/entity/balance"
	"github.com/anoriar/gophermart/internal/gophermart/repository/balance"
	"go.uber.org/zap"
)

type BalanceService struct {
	balanceRepository balance.BalanceRepositoryInterface
	logger            *zap.Logger
}

func NewBalanceService(balanceRepository balance.BalanceRepositoryInterface, logger *zap.Logger) *BalanceService {
	return &BalanceService{balanceRepository: balanceRepository, logger: logger}
}
func (service BalanceService) GetUserBalance(ctx context.Context, userID string) (balanceResponsePkg.BalanceResponseDto, error) {
	balanceResult, err := service.balanceRepository.GetBalanceByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			//не ошибка - если у пользователя отсутсвует запись о балансе - возвращаем 0
			return balanceResponsePkg.BalanceResponseDto{}, nil
		}
		service.logger.Error(err.Error())
		return balanceResponsePkg.BalanceResponseDto{}, err
	}
	return balanceResponsePkg.BalanceResponseDto{
		Current:   balanceResult.Balance,
		Withdrawn: balanceResult.Withdrawal,
	}, nil
}

func (service BalanceService) UpdateUserBalance(ctx context.Context, userID string, addSum float64, withdrawSum float64) error {
	err := service.balanceRepository.UpsertBalance(ctx, userID, func(currentBalance *balanceEntityPkg.Balance) balanceDtoPkg.UpdateBalanceDto {
		if currentBalance == nil {
			return balanceDtoPkg.UpdateBalanceDto{
				UserID:     userID,
				Balance:    addSum,
				Withdrawal: withdrawSum,
			}
		} else {
			return balanceDtoPkg.UpdateBalanceDto{
				UserID:     userID,
				Balance:    currentBalance.Balance + addSum - withdrawSum,
				Withdrawal: currentBalance.Withdrawal + withdrawSum,
			}
		}
	})
	if err != nil {
		service.logger.Error(err.Error())
		return err
	}

	service.logger.Info(fmt.Sprintf("user_id %s balance updated successfully", userID),
		zap.Float64("addSum", addSum),
		zap.Float64("withdrawSum", withdrawSum),
	)
	return nil
}
