package task

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/wayofthepie/task-store/pkg/model"
	"github.com/satori/go.uuid"
)

// TaskQueueKey : the key for a task queue
const taskQueueKey = "taskQ"

// TaskIDPrefix : the prefix for task id's
const taskIDPrefix = "task"

// Queue : a Queue allows pushing popping and reading
// of task information from a queue
type Queue interface {
	Push(spec *model.Spec) (*uuid.UUID, error)
}

// NewQueueImpl : build a QueueImpl
func NewQueueImpl(redis *redis.Client) *QueueImpl {
	return &QueueImpl{redis: redis}
}

// QueueImpl : redis implementation of a Queue.
type QueueImpl struct {
	redis *redis.Client
}

// Push : push the given TaskSpec on the queue
func (q *QueueImpl) Push(spec *model.Spec) (*uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("something went wrong when attempting to generate a new uuid: %s", err)
	}

	err = q.createTask(&id, *spec)
	if err != nil {
		return nil, err
	}

	_, err = q.pushOntoTaskQ(&id)

	return &id, err
}

func (q *QueueImpl) createTask(id *uuid.UUID, spec model.Spec) error {
	spec.ID = id
	applied, err := q.redis.SetNX(buildTaskId(id), &spec, 0).Result()
	if err != nil {
		return fmt.Errorf("creating model failed with error: %s", err.Error())
	}
	if !applied {
		return fmt.Errorf("task with id %s already exists: %s", id.String(), err.Error())
	}
	return nil
}

func (q *QueueImpl) pushOntoTaskQ(id *uuid.UUID) (int64, error) {
	length, err := q.redis.LPush(taskQueueKey, id.String()).Result()
	if err == nil {
		return length, nil
	}
	return 0, err
}

func buildTaskId(id *uuid.UUID) string {
	return fmt.Sprintf("%s:%s", taskIDPrefix, id.String())
}
