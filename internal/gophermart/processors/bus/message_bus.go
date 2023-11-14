package bus

import (
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
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
	switch msg.(type) {
	case message.OrderProcessMessage:
		bus.OrderProcessChan <- msg.(message.OrderProcessMessage)
	default:
		return domain_errors.ErrNotValidTypeOfMessage
	}
	return nil
}
