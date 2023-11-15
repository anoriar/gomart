package balance

import (
	"context"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	balanceResponsePkg "github.com/anoriar/gophermart/internal/gophermart/dto/responses/balance"
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
		if errors.Is(err, domain_errors.ErrNotFound) {
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

func (service BalanceService) SyncUserBalance(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}
