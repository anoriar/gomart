package balance

import (
	"context"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/requests"
	balanceResponsePkg "github.com/anoriar/gophermart/internal/gophermart/balance/dto/responses"
	"github.com/anoriar/gophermart/internal/gophermart/balance/repository/balance"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/shared/errors"
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
		if errors.Is(err, errors2.ErrNotFound) {
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

func (service BalanceService) AddUserBalance(ctx context.Context, userID string, sum float64) error {
	err := service.balanceRepository.AddUserBalance(ctx, userID, sum)
	if err != nil {
		service.logger.Error(err.Error())
		return err
	}
	return nil
}

func (service BalanceService) WithdrawUserBalance(ctx context.Context, userID string, withdrawDto requests.WithdrawDto) error {
	err := service.balanceRepository.WithdrawUserBalance(ctx, userID, withdrawDto)
	if err != nil {
		service.logger.Error(err.Error())
		return err
	}
	return nil
}
