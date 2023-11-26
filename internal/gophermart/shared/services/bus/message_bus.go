package bus

import (
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/order/processors/orderprocess/message"
	"github.com/anoriar/gophermart/internal/gophermart/shared/errors"
	"go.uber.org/zap"
)

type MessageBus struct {
	OrderProcessChan chan message.OrderProcessMessage
	logger           *zap.Logger
}

func NewMessageBus(logger *zap.Logger) *MessageBus {
	return &MessageBus{
		OrderProcessChan: make(chan message.OrderProcessMessage, 100),
		logger:           logger,
	}
}

func (bus *MessageBus) SendMessage(msg interface{}) error {
	switch m := msg.(type) {
	case message.OrderProcessMessage:
		bus.OrderProcessChan <- m
	default:
		return errors.ErrNotValidTypeOfMessage
	}
	bus.logger.Info(fmt.Sprintf("Message sended: %s", msg))
	return nil
}
