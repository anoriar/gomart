package fetcher

import (
	"context"
	accrualPkg "github.com/anoriar/gophermart/internal/gophermart/order/dto/accrual"
	"github.com/anoriar/gophermart/internal/gophermart/order/entity"
	"github.com/anoriar/gophermart/internal/gophermart/order/repository/accrual"
	"github.com/opentracing/opentracing-go"
)

type OrderFetchService struct {
	accrualRepository accrual.AccrualRepositoryInterface
}

func NewOrderFetchService(accrualRepository accrual.AccrualRepositoryInterface) *OrderFetchService {
	return &OrderFetchService{accrualRepository: accrualRepository}
}

func (service OrderFetchService) Fetch(ctx context.Context, orderEntity entity.Order) (entity.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderFetchService::Fetch")
	span.SetTag("orderID", orderEntity.ID)
	defer span.Finish()

	var status string
	accrual := orderEntity.Accrual

	extOrder, extOrderExists, err := service.accrualRepository.GetOrder(ctx, orderEntity.ID)
	if err != nil {
		return orderEntity, err
	}
	if !extOrderExists {
		status = entity.ProcessingStatus
	} else {
		switch extOrder.Status {
		case accrualPkg.AccrualProcessedStatus:
			status = entity.ProcessedStatus
			accrual = extOrder.Accrual
		case accrualPkg.AccrualInvalidStatus:
			status = entity.InvalidStatus
		default:
			status = entity.ProcessingStatus
		}
	}

	orderEntity.Status = status
	orderEntity.Accrual = accrual

	return orderEntity, nil
}
