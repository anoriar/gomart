package orderprocess

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/order/dto/repository"
	"github.com/anoriar/gophermart/internal/gophermart/order/entity"
	"github.com/anoriar/gophermart/internal/gophermart/order/processors/orderprocess/message"
	"github.com/anoriar/gophermart/internal/gophermart/order/services"
	repositoryDtoPkg "github.com/anoriar/gophermart/internal/gophermart/shared/dto/repository"
	"github.com/anoriar/gophermart/internal/gophermart/shared/services/bus"
	"go.uber.org/zap"
	"math"
	"time"
)

const (
	batchSize = 100
)

type OrderSyncFailedProcessor struct {
	orderService services.OrderServiceInterface
	logger       *zap.Logger
	msgBus       bus.MessageBusInterface
}

func NewOrderSyncFailedProcessor(
	ctx context.Context,
	orderService services.OrderServiceInterface,
	logger *zap.Logger,
	msgBus bus.MessageBusInterface,
) *OrderSyncFailedProcessor {
	instance := &OrderSyncFailedProcessor{
		orderService: orderService,
		logger:       logger,
		msgBus:       msgBus,
	}

	go instance.flush(ctx)

	return instance
}

func (p *OrderSyncFailedProcessor) flush(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	orderFilter := repository.OrdersFilterDto{
		Statuses: []string{entity.ProcessingStatus, entity.NewStatus},
	}
	for {
		select {
		case <-ticker.C:
			total, err := p.orderService.GetOrdersTotal(ctx, orderFilter)
			if err != nil {
				p.logger.Error(err.Error())
				continue
			}
			stepCount := int(math.Ceil(float64(total) / float64(batchSize)))
			for i := 0; i < stepCount; i++ {
				orders, err := p.orderService.GetOrders(ctx, repository.OrdersQuery{
					Filter: orderFilter,
					Pagination: repositoryDtoPkg.PaginationDto{
						Limit:  batchSize,
						Offset: i * batchSize,
					},
				})
				if err != nil {
					p.logger.Error(err.Error())
					continue
				}

				for _, order := range orders {
					err := p.msgBus.SendMessage(message.OrderProcessMessage{OrderID: order.ID})
					if err != nil {
						p.logger.Error(err.Error())
						continue
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
