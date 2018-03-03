package task

import (
	"net/http"
	"github.com/wayofthepie/jobby-taskman/pkg/model"
	"io/ioutil"
	"fmt"
)

type Handler interface {
	CreateTaskHandler(w http.ResponseWriter, r *http.Request)
}

type HandlerImpl struct {
	taskQueue TaskQueue
}

func Test() {

}
func NewHandlerImpl(taskQueue TaskQueue) *HandlerImpl {
	return &HandlerImpl{taskQueue: taskQueue}
}

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

	taskId, err := h.taskQueue.Push(taskSpec)

	if err = taskSpec.UnmarshalBinary(body); err != nil {
		build500Error("failed to store task", err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(taskId))
}

func build500Error(customMsg string, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	errMsg := fmt.Sprintf("%s : %s", customMsg, err.Error())
	w.Write([]byte(errMsg))
	return
}
