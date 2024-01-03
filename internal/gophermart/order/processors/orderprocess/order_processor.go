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
	orderService services.OrderServiceInterface,
	logger *zap.Logger,
	msgChan chan message.OrderProcessMessage,
) *OrderProcessor {
	instance := &OrderProcessor{
		orderService: orderService,
		logger:       logger,
		msgChan:      msgChan,
	}

	go instance.process()

	return instance
}

func (p *OrderProcessor) process() {
	for msg := range p.msgChan {
		err := p.orderService.ProcessOrder(context.Background(), msg.OrderID)
		if err != nil {
			p.logger.Error("process order error", zap.String("error", err.Error()))
		}
		continue
	}
}
