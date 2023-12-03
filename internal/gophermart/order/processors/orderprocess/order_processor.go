package orderprocess

import (
	"context"
	"github.com/anoriar/gophermart/internal/gophermart/order/processors/orderprocess/message"
	"github.com/anoriar/gophermart/internal/gophermart/order/services"
	"go.uber.org/zap"
)

type OrderProcessor struct {
	orderService services.OrderServiceInterface
	logger       *zap.Logger
	msgChan      chan message.OrderProcessMessage
}

func NewOrderProcessor(
	ctx context.Context,
	orderService services.OrderServiceInterface,
	logger *zap.Logger,
	msgChan chan message.OrderProcessMessage,
) *OrderProcessor {
	instance := &OrderProcessor{
		orderService: orderService,
		logger:       logger,
		msgChan:      msgChan,
	}

	go instance.process(ctx)

	return instance
}

func (p *OrderProcessor) process(ctx context.Context) {
	for msg := range p.msgChan {
		err := p.orderService.ProcessOrder(ctx, msg.OrderID)
		if err != nil {
			p.logger.Error("process order error", zap.String("error", err.Error()))
		}
		continue
	}
}
