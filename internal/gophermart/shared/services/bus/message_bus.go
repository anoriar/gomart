package bus

import (
	"github.com/anoriar/gophermart/internal/gophermart/order/processors/message"
	"github.com/anoriar/gophermart/internal/gophermart/shared/errors"
)

type MessageBus struct {
	OrderProcessChan chan message.OrderProcessMessage
}

func NewMessageBus() *MessageBus {
	return &MessageBus{
		OrderProcessChan: make(chan message.OrderProcessMessage, 100),
	}
}

func (bus *MessageBus) SendMessage(msg interface{}) error {
	switch m := msg.(type) {
	case message.OrderProcessMessage:
		bus.OrderProcessChan <- m
	default:
		return errors.ErrNotValidTypeOfMessage
	}
	return nil
}
