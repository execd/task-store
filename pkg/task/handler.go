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
	taskQueue    Queue
	eventService EventService
}

// NewHandlerImpl creates a new HandlerImpl
func NewHandlerImpl(taskQueue Queue, eventService EventService) *HandlerImpl {
	return &HandlerImpl{taskQueue: taskQueue, eventService: eventService}
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

	if err := h.eventService.PublishWork(taskSpec); err != nil {
		// TODO Cleanup created task as it will never be used
		// it may also be the case that publish fails due to limits on the queue
		build500Error("failed to send event for created task", err, w)
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
