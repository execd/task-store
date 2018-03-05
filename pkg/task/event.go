package task

import (
	"github.com/NeowayLabs/wabbit"
	"github.com/wayofthepie/task-store/pkg/model"
)

const queueName = "work_queue"

// EventService : interface for an EventService
type EventService interface {
	PublishWork(spec *model.TaskSpec) error
}

// AmqpEventService : amqp implementation of an event service
type AmqpEventService struct {
	channel wabbit.Channel
	queue   wabbit.Queue
}

// NewAmqpEventService : build a new AmqpEventService with the given Channel
func NewAmqpEventService(channel wabbit.Channel) (*AmqpEventService, error) {
	q, err := channel.QueueDeclare(
		queueName,
		wabbit.Option{
			"durable":    true,
			"autoDelete": false,
			"exclusive":  false,
			"noWait":     false,
		},
	)

	if err != nil {
		return nil, err
	}

	return &AmqpEventService{channel: channel, queue: q}, nil
}

// PublishWork : publish the given TaskSpec on the work queue
func (s *AmqpEventService) PublishWork(spec *model.TaskSpec) error {
	body, err := spec.MarshalBinary()
	if err != nil {
		return err
	}

	s.channel.Publish(
		"",             // exchange
		s.queue.Name(), // routing key
		body,           // mandatory
		wabbit.Option{
			"deliveryMode": 2,
			"contentType":  "application/json",
		})
	return nil
}
