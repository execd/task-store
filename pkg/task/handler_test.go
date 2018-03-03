package task

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/wayofthepie/task-store/mocks"
	"github.com/wayofthepie/task-store/pkg/model"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Handler", func() {
	Describe("CreateTaskHandler", func() {

		taskQueueMock := new(mocks.TaskQueue)
		handler := NewHandlerImpl(taskQueueMock)

		It("should push a task on the queue", func() {
			// Arrange
			expectedtaskID := "task:1"
			json := `{"name": "test", "image": "alpine", "init": "init.sh"}`
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader([]byte(json)))
			taskSpec := &model.TaskSpec{Name: "test", Image: "alpine", Init: "init.sh"}
			writer := httptest.NewRecorder()
			taskQueueMock.On("Push", taskSpec).Return(expectedtaskID, nil)

			// Act
			handler.CreateTaskHandler(writer, req)

			// Assert
			assert.Equal(context, 201, writer.Code)
			assert.Equal(context, expectedtaskID, writer.Body.String())
		})
	})
})
