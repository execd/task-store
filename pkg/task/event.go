package task

import (
	"encoding/json"
	"fmt"
	"github.com/execd/task-store/pkg/model"
	"github.com/execd/task-store/pkg/rabbit"
)

// EventManager : interface for an event listener
type EventManager interface {
	PublishWork(task *model.Spec) error
	ListenForProgress(quit <-chan int) (<-chan model.Status, <-chan error)
}

// EventManagerImpl : implementation of an event listener
type EventManagerImpl struct {
	rabbit rabbit.Service
}

// NewEventManagerImpl : build a ListenerImpl
func NewEventManagerImpl(rabbit rabbit.Service) (*EventManagerImpl, error) {
	return &EventManagerImpl{rabbit: rabbit}, nil
}

// PublishWork : publish a task
func (e *EventManagerImpl) PublishWork(task *model.Spec) error {
	fmt.Printf("Publishing work for task %s\n", task.ID.String())
	return e.rabbit.PublishWork(task)
}

// ListenForProgress : listen for task progress
func (e *EventManagerImpl) ListenForProgress(quit <-chan int) (<-chan model.Status, <-chan error) {
	status := make(chan model.Status, 100)
	errors := make(chan error)
	incoming := e.rabbit.GetTaskStatusQueueChan()
	go func() {
		defer close(status)
		for {
			select {
			case msg, ok := <-incoming:
				if !ok {
					fmt.Println("Incoming task complete event channel not ok!")
					continue
				}
				info := new(model.Status)
				err := json.Unmarshal(msg.Body(), info)
				if err != nil {
					fmt.Println("error occurred unmarshalling data")
					errors <- fmt.Errorf("error occurred unmarshalling data (%s) : %s", string(msg.Body()[:]), err.Error())
				} else {
					i := *info
					status <- i
				}
				msg.Ack(false)
			case <-quit:
				fmt.Println("Stopping task listener.")
				return
			}
		}
	}()
	return status, errors
}
