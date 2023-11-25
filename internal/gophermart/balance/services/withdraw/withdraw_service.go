package withdraw

import (
	"context"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/balance/dto/requests"
	withdrawalPkg "github.com/anoriar/gophermart/internal/gophermart/balance/entity"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/balance/errors"
	"github.com/anoriar/gophermart/internal/gophermart/balance/repository/withdrawal"
	balanceServicePkg "github.com/anoriar/gophermart/internal/gophermart/balance/services/balance"
	"github.com/anoriar/gophermart/internal/gophermart/balance/services/withdraw/internal/factory"
	"github.com/anoriar/gophermart/internal/gophermart/order/errors"
	errors3 "github.com/anoriar/gophermart/internal/gophermart/shared/errors"
	"github.com/anoriar/gophermart/internal/gophermart/shared/services/validator/idvalidator"
	"go.uber.org/zap"
)

type WithdrawService struct {
	withdrawalRepository withdrawal.WithdrawalRepositoryInterface
	balanceService       balanceServicePkg.BalanceServiceInterface
	withdrawalFactory    *factory.WithdrawalFactory
	idValidator          idvalidator.IDValidatorInterface
	logger               *zap.Logger
}

func NewWithdrawService(
	withdrawalRepository withdrawal.WithdrawalRepositoryInterface,
	balanceService balanceServicePkg.BalanceServiceInterface,
	idValidator idvalidator.IDValidatorInterface,
	logger *zap.Logger,
) *WithdrawService {
	return &WithdrawService{
		withdrawalRepository: withdrawalRepository,
		balanceService:       balanceService,
		withdrawalFactory:    factory.NewWithdrawalFactory(),
		idValidator:          idValidator,
		logger:               logger,
	}
}

func (service WithdrawService) Withdraw(ctx context.Context, userID string, withdrawDto requests.WithdrawDto) error {
	orderIDValid := service.idValidator.Validate(withdrawDto.Order)
	if !orderIDValid {
		return errors.ErrOrderNumberNotValid
	}
	currentBalance, err := service.balanceService.GetUserBalance(ctx, userID)
	if err != nil {
		service.logger.Error(err.Error())
		return fmt.Errorf("%w: %v", errors3.ErrInternalError, err)
	}
	//недостаточно средств
	if currentBalance.Current < withdrawDto.Sum {
		return errors2.ErrInsufficientFunds
	}

	withdrawalEntity := service.withdrawalFactory.CreateEntityFromRequest(userID, withdrawDto)
	err = service.withdrawalRepository.CreateWithdrawal(ctx, withdrawalEntity)
	if err != nil {
		service.logger.Error(err.Error())
		return fmt.Errorf("%w: %v", errors3.ErrInternalError, err)
	}

	err = service.balanceService.UpdateUserBalance(ctx, userID, 0, withdrawDto.Sum)
	if err != nil {
		service.logger.Error(err.Error())
		return fmt.Errorf("%w: %v", errors3.ErrInternalError, err)
	}

	return nil
}

func (service WithdrawService) GetWithdrawalsByUserID(ctx context.Context, userID string) ([]withdrawalPkg.Withdrawal, error) {
	result, err := service.withdrawalRepository.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		service.logger.Error(err.Error())
		return nil, fmt.Errorf("%w: %v", errors3.ErrInternalError, err)
	}
	return result, nil
}
