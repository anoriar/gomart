package accrual

import (
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	accrualPkg "github.com/anoriar/gophermart/internal/gophermart/dto/accrual"
	"github.com/anoriar/gophermart/internal/gophermart/repository/accrual"
	"go.uber.org/zap"
)

type AccrualService struct {
	accrualRepository accrual.AccrualRepositoryInterface
	logger            *zap.Logger
}

func NewAccrualService(accrualRepository accrual.AccrualRepositoryInterface, logger *zap.Logger) *AccrualService {
	return &AccrualService{accrualRepository: accrualRepository, logger: logger}
}

func (service AccrualService) GetOrder(orderId string) (accrualPkg.AccrualOrderDto, error) {
	order, err := service.accrualRepository.GetOrder(orderId)
	if err != nil {
		if !errors.Is(err, domain_errors.ErrNotFound) {
			service.logger.Error(err.Error())
		}

		return accrualPkg.AccrualOrderDto{}, err
	}
	return order, nil
}
