package event

import (
	"fmt"
	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
	"github.com/wayofthepie/task-store/pkg/model"
	"log"
)

const queueName = "work_queue"

// Rabbit : interface for building rabbit services
type Rabbit interface {
	GetTaskStatusQueueChan() <-chan wabbit.Delivery
	PublishWork(task *model.Spec) error
}

// RabbitImpl : service to interact with rabbit
type RabbitImpl struct {
	connection          wabbit.Conn
	channel             wabbit.Channel
	taskStatusQueue     wabbit.Queue
	taskStatusQueueName string
	taskStatusQueueChan <-chan wabbit.Delivery
	workQueueName       string
}

// NewRabbitServiceImpl : build a new connection to rabbitmq
func NewRabbitServiceImpl(address string) (*RabbitImpl, error) {

	r := &RabbitImpl{}
	r.initialize(address)

	return r, nil
}

// GetTaskStatusQueueChan : get the work queue channel
func (r *RabbitImpl) GetTaskStatusQueueChan() <-chan wabbit.Delivery {
	return r.taskStatusQueueChan
}

// PublishWork : publish work on the work queue
func (r *RabbitImpl) PublishWork(task *model.Spec) error {
	data, err := task.MarshalBinary()
	if err != nil {
		return err
	}
	opts := wabbit.Option{
		"contentType": "application/json",
	}
	return r.channel.Publish("", r.workQueueName, data, opts)
}

func (r *RabbitImpl) initialize(address string) {
	c := make(chan wabbit.Error)
	fmt.Println("Initializing")
	go func() {
		err := <-c
		log.Println("reconnect: ", err.Error())
		r.initialize(address)
	}()

	conn, err := amqp.Dial(address)
	if err != nil {
		panic("cannot connect")
	}
	conn.NotifyClose(c)

	ch, err := conn.Channel()
	if err != nil {
		panic("cannot create channel")
	}

	r.connection = conn
	r.channel = ch
	r.initializeWorkQueueConsumer()
	r.declareTaskQueue()
}

func (r *RabbitImpl) initializeWorkQueueConsumer() {
	workQueue, err := r.channel.QueueDeclare(
		queueName,
		wabbit.Option{
			"durable":    true,
			"autoDelete": false,
			"exclusive":  false,
			"noWait":     false,
		},
	)
	r.channel.Qos(
		5,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic("Could not setup work_queue")
	}

	r.workQueueName = workQueue.Name()
}

func (r *RabbitImpl) declareTaskQueue() {
	name := "task_status_queue"
	taskStatusQueue, err := r.channel.QueueDeclare(
		name,
		wabbit.Option{
			"durable":    true,
			"autoDelete": false,
			"exclusive":  false,
			"noWait":     false,
		},
	)
	if err != nil {
		panic("Could not setup task_queue")
	}
	taskStatusQueueChan, err := r.channel.Consume(
		taskStatusQueue.Name(),
		"",
		wabbit.Option{
			"auto-ack":  false,
			"exclusive": false,
			"no-local":  false,
			"no-wait":   false,
		},
	)
	if err != nil {
		panic("Could not setup task_queue")
	}
	r.taskStatusQueueChan = taskStatusQueueChan
	r.taskStatusQueue = taskStatusQueue
	r.taskStatusQueueName = name
}
