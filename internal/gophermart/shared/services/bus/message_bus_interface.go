package bus

type MessageBusInterface interface {
	SendMessage(msg interface{}) error
}
