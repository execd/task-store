package task

import (
	"fmt"
	"github.com/execd/task-store/pkg/model"
	"github.com/execd/task-store/pkg/util"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
)

const taskQueueName = "taskQ"
const executingQueueName = "executing"
const taskPrefix = "task"
const taskCreatedChannel = "taskCreated"

// Store : a Store allows pushing popping and reading
// of task information from a queue
type Store interface {
	StoreTask(task model.Spec) (*uuid.UUID, error)
	GetTask(id *uuid.UUID) (*model.Spec, error)
	PushTask(id *uuid.UUID) (*uuid.UUID, error)
	PopTask() (*uuid.UUID, error)
	AddTaskToExecutingSet(id *uuid.UUID) error
	RemoveTaskFromExecutingSet(id *uuid.UUID) error
	IsTaskExecuting(id *uuid.UUID) (bool, error)
	PublishTaskCreatedEvent(id *uuid.UUID) error
	ListenForTaskCreatedEvents() <-chan uuid.UUID
}

// NewStoreImpl : build a StoreImpl
func NewStoreImpl(redis *redis.Client, uuidGen util.UUIDGen) *StoreImpl {
	return &StoreImpl{redis: redis, uuidGen: uuidGen}
}

// StoreImpl : redis implementation of a Store.
type StoreImpl struct {
	redis   *redis.Client
	uuidGen util.UUIDGen
}

// StoreTask : store the given task
func (s *StoreImpl) StoreTask(task model.Spec) (*uuid.UUID, error) {
	id, err := s.uuidGen.GenV4()
	if err != nil {
		return nil, err
	}
	task.ID = &id
	created, err := s.redis.SetNX(buildTaskKey(&id), &task, 0).Result()
	if err != nil {
		return nil, fmt.Errorf("storing task with id %s failed", task.ID)
	}
	if !created {
		return nil, fmt.Errorf("task with id %s already exists", task.ID.String())
	}
	return task.ID, nil
}

// GetTask : retrieve the task with the given id
func (s *StoreImpl) GetTask(id *uuid.UUID) (*model.Spec, error) {
	task, err := s.redis.Get(buildTaskKey(id)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve task with id %s", id.String())
	}

	taskSpec := new(model.Spec)
	if err := taskSpec.UnmarshalBinary([]byte(task)); err != nil {
		return nil, fmt.Errorf("failed to build task with id %s from retrieved data %s", id.String(), task)
	}
	return taskSpec, nil
}

// PushTask : push the given TaskSpec on the task queue
func (s *StoreImpl) PushTask(id *uuid.UUID) (*uuid.UUID, error) {
	_, err := s.redis.LPush(taskQueueName, id.String()).Result()
	if err != nil {
		return nil, err
	}
	return id, nil
}

// PopTask : get the next task
func (s *StoreImpl) PopTask() (*uuid.UUID, error) {
	results, err := s.redis.BRPop(0, taskQueueName).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve next task to execute : %s", err.Error())
	}

	stringID := results[1]

	id := new(uuid.UUID)
	_ = id.UnmarshalText([]byte(stringID))

	return id, nil
}

// AddTaskToExecutingSet : move a task to the executing set
func (s *StoreImpl) AddTaskToExecutingSet(id *uuid.UUID) error {
	_, err := s.redis.SAdd(executingQueueName, id.String()).Result()
	if err != nil {
		return fmt.Errorf("failed to add task to executing set : %s", err.Error())
	}
	return nil
}

// RemoveTaskFromExecutingSet : remove task from the executing set
func (s *StoreImpl) RemoveTaskFromExecutingSet(id *uuid.UUID) error {
	i, err := s.redis.SRem(executingQueueName, id.String()).Result()
	println(i)
	if err != nil {
		return fmt.Errorf("failed to remove task %s : %s", id.String(), err.Error())
	}
	return nil
}

// IsTaskExecuting : true if a task is executing, false otherwise
func (s *StoreImpl) IsTaskExecuting(id *uuid.UUID) (bool, error) {
	return s.redis.SIsMember(executingQueueName, id.String()).Result()
}

// PublishTaskCreatedEvent : publish a task created event to the
// task created redis channel
func (s *StoreImpl) PublishTaskCreatedEvent(id *uuid.UUID) error {
	_, err := s.redis.Publish(taskCreatedChannel, id.String()).Result()
	if err != nil {
		return fmt.Errorf("failed to publish task created event : %s", err.Error())
	}
	return nil
}

// ListenForTaskCreatedEvents : get a channel where task
// created events will be pushed
// TODO : How to test this?
func (s *StoreImpl) ListenForTaskCreatedEvents() <-chan uuid.UUID {
	ids := make(chan uuid.UUID)
	sub := s.redis.Subscribe(taskCreatedChannel)
	go func() {
		for msg := range sub.Channel() {
			fmt.Printf("Got event for %s", msg.Payload)
			id, err := uuid.FromString(msg.Payload)
			if err != nil {
				fmt.Printf("failed to build id from %s\n", msg.Payload)
			} else {
				ids <- id
			}
		}
	}()
	return ids
}

func buildTaskKey(id *uuid.UUID) string {
	return fmt.Sprintf("%s:%s", taskPrefix, id.String())
}
