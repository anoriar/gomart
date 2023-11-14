package fetcher

import (
	accrualPkg "github.com/anoriar/gophermart/internal/gophermart/dto/accrual"
	"github.com/anoriar/gophermart/internal/gophermart/entity/order"
	"github.com/anoriar/gophermart/internal/gophermart/repository/accrual"
)

type OrderFetchService struct {
	accrualRepository accrual.AccrualRepositoryInterface
}

func NewOrderFetchService(accrualRepository accrual.AccrualRepositoryInterface) *OrderFetchService {
	return &OrderFetchService{accrualRepository: accrualRepository}
}

func (service OrderFetchService) Fetch(orderEntity order.Order) (order.Order, error) {
	status := orderEntity.Status
	accrual := orderEntity.Accrual

	extOrder, extOrderExists, err := service.accrualRepository.GetOrder(orderEntity.Id)
	if err != nil {
		return orderEntity, err
	}
	if extOrderExists == false {
		status = order.InvalidStatus
	} else {
		switch extOrder.Status {
		case accrualPkg.AccrualProcessedStatus:
			status = order.ProcessedStatus
			accrual = extOrder.Accrual
		case accrualPkg.AccrualInvalidStatus:
			status = order.InvalidStatus
		default:
			status = order.ProcessingStatus
		}
	}

	orderEntity.Status = status
	orderEntity.Accrual = accrual

	return orderEntity, nil
}
