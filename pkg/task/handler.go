package task

import (
	"fmt"
	"github.com/wayofthepie/task-store/pkg/model"
	"io/ioutil"
	"net/http"
)

// Handler : interface for a task Handler
type Handler interface {
	CreateTaskHandler(w http.ResponseWriter, r *http.Request)
}

// HandlerImpl : implementation of a task Handler
type HandlerImpl struct {
	taskQueue Queue
}

// NewHandlerImpl creates a new HandlerImpl
func NewHandlerImpl(taskQueue Queue) *HandlerImpl {
	return &HandlerImpl{taskQueue: taskQueue}
}

// CreateTaskHandler handles task creation requests
func (h *HandlerImpl) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	taskSpec := new(model.TaskSpec)

	if err = taskSpec.UnmarshalBinary(body); err != nil {
		build500Error("failed to parse request body", err, w)
		return
	}

	taskID, err := h.taskQueue.Push(taskSpec)
	if err != nil {
		build500Error("failed to store task", err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(taskID))
}

func build500Error(customMsg string, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	errMsg := fmt.Sprintf("%s : %s", customMsg, err.Error())
	w.Write([]byte(errMsg))
}
