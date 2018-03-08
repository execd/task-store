package task

import (
	"github.com/go-redis/redis"
	"github.com/wayofthepie/task-store/pkg/model"
	"github.com/satori/go.uuid"
)

// TaskQueueKey : the key for a task queue
const taskQueueKey = "taskQ"

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
func (q *QueueImpl) Push(id *uuid.UUID) (*uuid.UUID, error) {
	_, err := q.redis.LPush(taskQueueKey, id.String()).Result()
	if err != nil {
		return nil, err
	}
	return id, nil
}
