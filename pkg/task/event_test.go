package task_test

import (
	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/execd/task-store/mocks"
	"github.com/execd/task-store/pkg/model"
	"github.com/execd/task-store/pkg/task"
	. "github.com/onsi/ginkgo"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"time"
)

var _ = Describe("event", func() {
	Describe("listen for task progress", func() {
		var rabbitMock *mocks.Rabbit
		var eventListener *task.EventManagerImpl

		BeforeEach(func() {
			rabbitMock = &mocks.Rabbit{}
			eventListener, _ = task.NewEventManagerImpl(rabbitMock)
		})

		It("should quit when quit channel has item", func() {
			// Arrange
			quit := make(chan int)
			defer close(quit)
			rabbitMock.On("GetTaskStatusQueueChan").Return(nil)
			timeout := time.After(time.Millisecond * 5)

			// Act
			status, _ := eventListener.ListenForProgress(quit)

			quit <- 1

			// Assert
			select {
			case _, ok := <-status:
				if !ok {
					status = nil
					break
				}
				assert.Fail(context, "Should not loop.")
			case <-timeout:
				assert.Fail(context, "Timed out waiting for channel to close")
			}
		})

		It("should handle message when channel has message", func() {
			// Arrange
			quit := make(chan int, 1)
			defer close(quit)
			id := uuid.Must(uuid.NewV4())
			expectedInfo := &model.Info{
				ID: &id,
			}
			data, err := expectedInfo.MarshalBinary()
			assert.Nil(context, err)

			rabbitMock.On("GetTaskStatusQueueChan").Return(buildMsgChan(data))
			timeout := time.After(time.Millisecond * 5)

			// Act
			status, _ := eventListener.ListenForProgress(quit)

			// Assert
			select {

			case info, ok := <-status:
				if !ok {
					status = nil
					quit <- 1
					assert.Fail(context, "Channel failed unexpectedly")
				}
				assert.Equal(context, *expectedInfo, info)
				break
			case <-timeout:
				assert.Fail(context, "Timed out waiting for channel to close, or data to be received")
			}
		})

		It("should add error to error channel if msg fails to be decoded", func() {
			// Arrange
			quit := make(chan int, 1)
			defer close(quit)

			rabbitMock.On("GetTaskStatusQueueChan").Return(buildMsgChan([]byte("not right")))
			timeout := time.After(time.Millisecond * 5)

			// Act
			status, errors := eventListener.ListenForProgress(quit)

			// Assert
			select {
			case <-status:
				assert.Fail(context, "Should not read a status.")
			case err := <-errors:
				assert.Contains(context, err.Error(), "error occurred unmarshalling data")
			case <-timeout:
				assert.Fail(context, "Timed out waiting for channel to close, or error to be received")
			}
		})
	})
})

func buildMsgChan(data []byte) <-chan wabbit.Delivery {
	msgs := make(chan wabbit.Delivery, 1)
	msg := server.NewDelivery(nil, data, 0, "", wabbit.Option{})
	msgs <- msg
	return msgs
}
