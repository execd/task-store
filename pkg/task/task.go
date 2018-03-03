package task

import (
	"github.com/go-redis/redis"
	"fmt"
	"github.com/wayofthepie/task-store/pkg/model"
)

const TaskQueueKey = "taskQ"
const TaskIdPrefix = "task"
const LastTaskId = "task:id"

type TaskQueue interface {
	GetTaskInfo(taskId string) (*model.TaskSpec, error)
	Push(spec *model.TaskSpec) (string, error)
}

type TaskQueueImpl struct {
	redis *redis.Client
}

func (q *TaskQueueImpl) GetTaskInfo(taskId string) (*model.TaskSpec, error) {
	scmd := q.redis.Get(taskId)
	taskData, err := scmd.Result()
	if err != nil {
		return nil, fmt.Errorf("retrieving model with id %s failed with: %s", taskId, err.Error())
	}
	taskSpec := new(model.TaskSpec)
	taskSpec.UnmarshalBinary([]byte(taskData))
	return taskSpec, nil
}

func (q *TaskQueueImpl) Push(spec *model.TaskSpec) (string, error) {
	num, err := q.getNextTaskNumber()
	if err != nil {
		return "", err
	}

	taskId := buildTaskId(num)

	err = q.createTask(taskId, spec)
	if err != nil {
		return "", err
	}

	_, err = q.pushOntoTaskQ(taskId)

	return taskId, err
}

func (q *TaskQueueImpl) getNextTaskNumber() (int64, error) {
	icmd := q.redis.Incr(LastTaskId)
	if id, err := icmd.Result(); err != nil {
		return 0, fmt.Errorf("reserving an id failed with: %s", err.Error())
	} else {
		return id, nil
	}
}

func (q *TaskQueueImpl) createTask(taskId string, spec *model.TaskSpec) error {
	icmd := q.redis.Set(taskId, spec, 0)
	if _, err := icmd.Result(); err != nil {
		return fmt.Errorf("creating model failed with error: %s", err.Error())
	}
	return nil
}

func (q *TaskQueueImpl) pushOntoTaskQ(taskId string) (int64, error) {
	icmd := q.redis.LPush(TaskQueueKey, taskId)
	if length, err := icmd.Result(); err == nil {
		return length, nil
	} else {
		return 0, err
	}
}

func buildTaskId(num int64) string {
	return fmt.Sprintf("%s:%d", TaskIdPrefix, num)
}
