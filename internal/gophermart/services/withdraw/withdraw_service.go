package withdraw

import (
	"context"
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/domainerrors"
	"github.com/anoriar/gophermart/internal/gophermart/dto/requests/withdraw"
	withdrawalPkg "github.com/anoriar/gophermart/internal/gophermart/entity/withdrawal"
	"github.com/anoriar/gophermart/internal/gophermart/repository/withdrawal"
	balanceServicePkg "github.com/anoriar/gophermart/internal/gophermart/services/balance"
	"github.com/anoriar/gophermart/internal/gophermart/services/validator/idvalidator"
	"github.com/anoriar/gophermart/internal/gophermart/services/withdraw/internal/factory"
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

func (service WithdrawService) Withdraw(ctx context.Context, userID string, withdrawDto withdraw.WithdrawDto) error {
	orderIDValid := service.idValidator.Validate(withdrawDto.Order)
	if !orderIDValid {
		return domainerrors.ErrOrderNumberNotValid
	}
	currentBalance, err := service.balanceService.GetUserBalance(ctx, userID)
	if err != nil {
		service.logger.Error(err.Error())
		return fmt.Errorf("%w: %v", domainerrors.ErrInternalError, err)
	}
	//недостаточно средств
	if currentBalance.Current < withdrawDto.Sum {
		return domainerrors.ErrInsufficientFunds
	}

	withdrawalEntity := service.withdrawalFactory.CreateEntityFromRequest(userID, withdrawDto)
	err = service.withdrawalRepository.CreateWithdrawal(ctx, withdrawalEntity)
	if err != nil {
		service.logger.Error(err.Error())
		return fmt.Errorf("%w: %v", domainerrors.ErrInternalError, err)
	}

	err = service.balanceService.UpdateUserBalance(ctx, userID, 0, withdrawDto.Sum)
	if err != nil {
		service.logger.Error(err.Error())
		return fmt.Errorf("%w: %v", domainerrors.ErrInternalError, err)
	}

	return nil
}

func (service WithdrawService) GetWithdrawalsByUserID(ctx context.Context, userID string) ([]withdrawalPkg.Withdrawal, error) {
	result, err := service.withdrawalRepository.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		service.logger.Error(err.Error())
		return nil, fmt.Errorf("%w: %v", domainerrors.ErrInternalError, err)
	}
	return result, nil
}
