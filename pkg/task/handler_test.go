package task

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/wayofthepie/task-store/mocks"
	"github.com/wayofthepie/task-store/pkg/model"
	"net/http"
	"net/http/httptest"
	"github.com/satori/go.uuid"
)

var _ = Describe("Handler", func() {
	Describe("CreateTaskHandler", func() {
		var writer *httptest.ResponseRecorder
		taskQueueMock := new(mocks.Queue)
		rabbitMock := new(mocks.Rabbit)
		taskSpec := &model.Spec{Name: "test", Image: "alpine", Init: "init.sh", InitArgs: []string{"10"}}
		handler := NewHandlerImpl(taskQueueMock, rabbitMock)
		expectedTaskID := uuid.Must(uuid.NewV4())
		json := `{"name": "test", "image": "alpine", "init": "init.sh", "initArgs":["10"]}`
		req, _ := http.NewRequest("POST", "/handle", bytes.NewReader([]byte(json)))

		BeforeEach(func() {
			writer = httptest.NewRecorder()
			taskQueueMock.On("Push", taskSpec).Return(&expectedTaskID, nil)
			rabbitMock.On("PublishWork", taskSpec).Return(nil)
		})

		It("should push a task on the queue", func() {
			// Act
			handler.CreateTaskHandler(writer, req)

			// Assert
			assert.Equal(context, 201, writer.Code)
			assert.Equal(context, expectedTaskID.String(), writer.Body.String())
		})

		It("should send an event if storing the task was successful", func() {
			// Act
			handler.CreateTaskHandler(writer, req)

			// Assert
			rabbitMock.AssertCalled(context, "PublishWork", taskSpec)
		})
	})
})
