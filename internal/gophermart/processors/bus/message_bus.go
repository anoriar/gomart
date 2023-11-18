package bus

import (
	"github.com/anoriar/gophermart/internal/gophermart/domainerrors"
	"github.com/anoriar/gophermart/internal/gophermart/processors/order/message"
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
		return domainerrors.ErrNotValidTypeOfMessage
	}
	return nil
}
