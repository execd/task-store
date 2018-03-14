package task

import (
	"fmt"
	"github.com/execd/task-store/pkg/event"
	"github.com/execd/task-store/pkg/model"
)

// EventManager : interface for an event listener
type EventManager interface {
	PublishWork(task *model.Spec) error
	ListenForProgress(quit <-chan int) <-chan model.Info
}

// EventManagerImpl : implementation of an event listener
type EventManagerImpl struct {
	rabbit event.Rabbit
}

// NewEventListenerImpl : build a ListenerImpl
func NewEventListenerImpl(rabbit event.Rabbit) (*EventManagerImpl, error) {
	return &EventManagerImpl{rabbit: rabbit}, nil
}

// PublishWork : publish a task
func (e *EventManagerImpl) PublishWork(task *model.Spec) error {
	return nil
}

// ListenForProgress : listen for task progress
func (e *EventManagerImpl) ListenForProgress(quit <-chan int) <-chan model.Info {
	status := make(chan model.Info)
	incoming := e.rabbit.GetTaskStatusQueueChan()
	go func() {
		defer close(status)
		for {
			select {
			case msg := <-incoming:
				info := new(model.Info)
				info.UnmarshalBinary(msg.Body())
				status <- *info
			case <-quit:
				fmt.Println("Stopping task listener.")
				return
			}
		}
	}()
	return status
}
