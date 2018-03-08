package event

import (
	"fmt"
	"github.com/NeowayLabs/wabbit"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"time"
)

// Listener : interface for execution of tasks
type Listener interface {
	ListenForTaskStatus() error
}

// ListenerImpl : default implementation of Listener
type ListenerImpl struct {
	rabbit Rabbit
}

// NewServiceImpl : build a ListenerImpl
func NewServiceImpl(rabbit Rabbit) (*ListenerImpl, error) {
	return &ListenerImpl{rabbit: rabbit}, nil
}

// ListenForTaskStatus : listen for and execute tasks
func (s *ListenerImpl) ListenForTaskStatus() error {
	fmt.Println("Listening for tasks")
	go func() {
		for {
			for msg := range s.rabbit.GetTaskStatusQueueChan() {
				go s.handleMsg(msg)
			}
			fmt.Println("Stopped listening for messages, waiting 2 seconds for reconnect...")
			time.Sleep(time.Second * 2)
		}
	}()
	return nil
}

func (s *ListenerImpl) handleMsg(msg wabbit.Delivery) {
	fmt.Println("received msg")
	taskInfo := new(task.Info)
	err := taskInfo.UnmarshalBinary(msg.Body())
	if err == nil {
		fmt.Println(string(msg.Body()))
	} else {
		fmt.Printf("received an error %s\n", err.Error())
	}
	msg.Ack(false)
}
