package task

import (
	. "github.com/onsi/ginkgo"
	"github.com/wayofthepie/jobby-taskman/mocks"
	"bytes"
	"net/http/httptest"
	"net/http"
	"github.com/stretchr/testify/assert"
	"github.com/wayofthepie/jobby-taskman/pkg/model"
)

var _ = Describe("Handler", func() {
	Describe("CreateTaskHandler", func() {

		taskQueueMock := new(mocks.TaskQueue)
		handler := NewHandlerImpl(taskQueueMock)

		It("should push a task on the queue", func() {
			// Arrange
			expectedTaskId := "task:1"
			json := `{"name": "test", "image": "alpine", "init": "init.sh"}`
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader([]byte(json)))
			taskSpec := &model.TaskSpec{Name:"test", Image:"alpine", Init:"init.sh"}
			writer := httptest.NewRecorder()
			taskQueueMock.On("Push", taskSpec).Return(expectedTaskId, nil)

			// Act
			handler.CreateTaskHandler(writer, req)

			// Assert
			assert.Equal(context, 201, writer.Code)
			assert.Equal(context, expectedTaskId, writer.Body.String())
		})
	})
})