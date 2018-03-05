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
		var writer *httptest.ResponseRecorder
		taskQueueMock := new(mocks.TaskQueue)
		eventServiceMock := new(mocks.EventService)
		taskSpec := &model.TaskSpec{Name: "test", Image: "alpine", Init: "init.sh", InitArgs: []string{"10"}}
		handler := NewHandlerImpl(taskQueueMock, eventServiceMock)
		expectedtaskID := "task:1"
		json := `{"name": "test", "image": "alpine", "init": "init.sh", "initArgs":["10"]}`
		req, _ := http.NewRequest("POST", "/handle", bytes.NewReader([]byte(json)))

		BeforeEach(func() {

			writer = httptest.NewRecorder()
			taskQueueMock.On("Push", taskSpec).Return(expectedtaskID, nil)
			eventServiceMock.On("PublishWork", taskSpec).Return(nil)
		})

		It("should push a task on the queue", func() {
			// Act
			handler.CreateTaskHandler(writer, req)

			// Assert
			assert.Equal(context, 201, writer.Code)
			assert.Equal(context, expectedtaskID, writer.Body.String())
		})

		It("should send an event if storing the task was successful", func() {
			// Act
			handler.CreateTaskHandler(writer, req)

			// Assert
			eventServiceMock.AssertCalled(context, "PublishWork", taskSpec)
		})
	})
})
