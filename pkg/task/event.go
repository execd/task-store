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
func (e *EventManagerImpl) ListenForProgress(quit <-chan int) (<-chan model.Info, <-chan error) {
	status := make(chan model.Info)
	errors := make(chan error)
	incoming := e.rabbit.GetTaskStatusQueueChan()
	go func() {
		defer close(status)
		for {
			select {
			case msg := <-incoming:
				info := new(model.Info)
				err := info.UnmarshalBinary(msg.Body())
				if err != nil {
					errors <-
						fmt.Errorf("error occurred unmarshalling data (%s) : %s", string(msg.Body()[:]), err.Error())
				} else {
					status <- *info
				}
			case <-quit:
				fmt.Println("Stopping task listener.")
				return
			}
		}
	}()
	return status, errors
}
