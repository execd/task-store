package route

import (
	"encoding/json"
	"fmt"
	"github.com/execd/task-store/pkg/model"
	"github.com/execd/task-store/pkg/task"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

// TaskHandler : interface for a task Handler
type TaskHandler interface {
	CreateTaskHandler(w http.ResponseWriter, r *http.Request)
}

// TaskHandlerImpl : implementation of a task Handler
type TaskHandlerImpl struct {
	taskStore task.Store
	config    *model.Config
}

// NewTaskHandlerImpl creates a new HandlerImpl
func NewTaskHandlerImpl(taskStore task.Store, config *model.Config) *TaskHandlerImpl {
	return &TaskHandlerImpl{taskStore: taskStore, config: config}
}

// CreateTask handles task creation requests
func (h *TaskHandlerImpl) CreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	taskSpec := new(model.Spec)
	err = taskSpec.UnmarshalBinary(body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	size, err := h.taskStore.TaskQueueSize()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	capacity := h.config.Manager.TaskQueueSize
	if size >= capacity {
		errStr := fmt.Sprintf("Failed to create task, task queue has reached its limit!")
		fmt.Println(errStr)
		http.Error(w, errStr, 500)
		return
	}

	id, err := h.taskStore.StoreTask(*taskSpec)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	queueSize, err := h.taskStore.PushTask(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Printf("Task %s added to task queue - remaining queue capacity is %d\n", id.String(), capacity-queueSize)

	h.taskStore.PublishTaskCreatedEvent(id)
	if err != nil {
		fmt.Printf("failed to publish task created event for %s: %s\n", id.String(), err.Error())
	}

	w.WriteHeader(201)
	w.Write([]byte(id.String()))
}

// GetTask : retrieve a task denoted by the given id
func (h *TaskHandlerImpl) GetTask(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	idStr := vars["id"]
	id, err := uuid.FromString(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to build id from %s : %s", idStr, err.Error()), 500)
		return
	}

	taskSpec, err := h.taskStore.GetTask(&id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data, err := json.Marshal(taskSpec)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}
