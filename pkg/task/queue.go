package task

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/wayofthepie/task-store/pkg/model"
)

// TaskQueueKey : the key for a task queue
const taskQueueKey = "taskQ"

// TaskIDPrefix : the prefix for task id's
const taskIDPrefix = "task"

// LastTaskID : the id of the key for the last task id
const lastTaskID = "task:id"

// Queue : a Queue allows pushing popping and reading
// of task information from a queue
type Queue interface {
	GetTaskInfo(taskID string) (*model.TaskSpec, error)
	Push(spec *model.TaskSpec) (string, error)
}

// NewQueueImpl : build a QueueImpl
func NewQueueImpl(redis *redis.Client) *QueueImpl {
	return &QueueImpl{redis: redis}
}

// QueueImpl : redis implementation of a Queue.
type QueueImpl struct {
	redis *redis.Client
}

// GetTaskInfo : get information on the task with the given id
func (q *QueueImpl) GetTaskInfo(taskID string) (*model.TaskSpec, error) {
	scmd := q.redis.Get(taskID)
	taskData, err := scmd.Result()
	if err != nil {
		return nil, fmt.Errorf("retrieving task with id %s failed with: %s", taskID, err.Error())
	}
	taskSpec := new(model.TaskSpec)
	taskSpec.UnmarshalBinary([]byte(taskData))
	return taskSpec, nil
}

// Push : push the given TaskSpec on the queue
func (q *QueueImpl) Push(spec *model.TaskSpec) (string, error) {
	num, err := q.getNextTaskNumber()
	if err != nil {
		return "", err
	}

	taskID := buildtaskID(num)

	err = q.createTask(taskID, spec)
	if err != nil {
		return "", err
	}

	_, err = q.pushOntoTaskQ(taskID)

	return taskID, err
}

func (q *QueueImpl) getNextTaskNumber() (int64, error) {
	id, err := q.redis.Incr(lastTaskID).Result()
	if err != nil {
		return 0, fmt.Errorf("reserving an id failed with: %s", err.Error())
	}
	return id, nil
}

func (q *QueueImpl) createTask(taskID string, spec *model.TaskSpec) error {
	_, err := q.redis.Set(taskID, spec, 0).Result()
	if err != nil {
		return fmt.Errorf("creating model failed with error: %s", err.Error())
	}
	return nil
}

func (q *QueueImpl) pushOntoTaskQ(taskID string) (int64, error) {
	length, err := q.redis.LPush(taskQueueKey, taskID).Result()
	if err == nil {
		return length, nil
	}
	return 0, err
}

func buildtaskID(num int64) string {
	return fmt.Sprintf("%s:%d", taskIDPrefix, num)
}
