package task

import (
	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/wayofthepie/task-store/pkg/model"
	"log"
)

var _ = Describe("AmqpEventService", func() {
	Describe("PublishWork", func() {
		It("should push the given spec on the work queue", func() {
			// Arrange
			amqpServer := "amqp://localhost:5670/%2f"
			expectedTaskSpec := &model.TaskSpec{Image: "alpine", Name: "test", Init: "init.sh"}

			fakeServer := server.NewServer(amqpServer)
			fakeServer.Start()

			mockConn, err := amqptest.Dial(amqpServer)
			failOnError(err)

			channel, err := mockConn.Channel()
			failOnError(err)

			service, err := NewAmqpEventService(channel)
			failOnError(err)

			msgs, err := channel.Consume(
				"work_queue", // queue
				"",           // consumer
				wabbit.Option{
					"auto-ack":  false, // auto-ack
					"exclusive": false, // exclusive
					"no-local":  false, // no-local
					"no-wait":   false, // no-wait
					"args":      nil,   // args
				},
			)
			failOnError(err)

			// Act
			service.PublishWork(expectedTaskSpec)

			// Assert
			for d := range msgs {
				taskSpec := new(model.TaskSpec)
				taskSpec.UnmarshalBinary(d.Body())
				assert.Equal(context, expectedTaskSpec, taskSpec)
				d.Ack(false)
				break
			}
		})
	})
})

func failOnError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
