package task

import (
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"github.com/wayofthepie/task-store/pkg/model"
	"fmt"
)

const taskPrefix = "task"

type Store interface {
	CreateTask(model.Spec) (*uuid.UUID, error)
	GetTask(uuid.UUID) (*model.Spec, error)
}

type StoreImpl struct {
	redis *redis.Client
}

func NewStoreImpl(client *redis.Client) *StoreImpl {
	return &StoreImpl{redis: client}
}

func (s *StoreImpl) CreateTask(task model.Spec) (*uuid.UUID, error) {
	created, err := s.redis.SetNX(buildTaskStoreId(*task.ID), &task, 0).Result()
	if err != nil {
		return nil, fmt.Errorf("storing task with id %s failed", task.ID)
	}
	if !created {
		return nil, fmt.Errorf("task with id %s already exists", task.ID.String())
	}
	return task.ID, nil
}

func (s *StoreImpl) GetTask(id uuid.UUID) (*model.Spec, error) {
	task, err := s.redis.Get(buildTaskStoreId(id)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve task with id %s", id.String())
	}

	taskSpec := new(model.Spec)
	if err := taskSpec.UnmarshalBinary([]byte(task)); err != nil {
		return nil, fmt.Errorf("failed to build task with id %s from retrieved data %s", id.String(), task)
	}
	return taskSpec, nil
}

func buildTaskStoreId(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", taskPrefix, id.String())
}
