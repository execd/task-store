package event

import "github.com/NeowayLabs/wabbit/amqp"

// NewRabbitConnection : build a new connection to rabbitmq
func NewRabbitConnection(address string) (*amqp.Conn, error) {
	return amqp.Dial(address)
}
