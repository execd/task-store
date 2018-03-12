package route

import (
	"bytes"
	"errors"
	"github.com/alicebob/miniredis"
	. "github.com/onsi/ginkgo"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/wayofthepie/task-store/pkg/route"
	"github.com/wayofthepie/task-store/mocks"
	"github.com/wayofthepie/task-store/pkg/store"
	"github.com/wayofthepie/task-store/pkg/task"
	"github.com/wayofthepie/task-store/pkg/uuidgen"
	"net/http"
	"net/http/httptest"
)

var context = GinkgoT()

var _ = Describe("task handler", func() {
	Describe("create task", func() {
		var taskStore *task.StoreImpl
		var directRedis *miniredis.Miniredis
		var handler *route.TaskHandlerImpl

		BeforeEach(func() {
			s, err := miniredis.Run()
			if err != nil {
				panic(err)
			}
			directRedis = s
			redis := store.NewClient(s.Addr())
			uuidGen := uuidgen.NewUUIDGenImpl()
			taskStore = task.NewStoreImpl(redis, uuidGen)
			handler = route.NewTaskHandlerImpl(taskStore)
		})

		It("should return an error if reading body fails", func() {
			// Arrange
			r := bytes.NewReader(nil)
			r.Seek(10, 0)
			req, _ := http.NewRequest("POST", "/handle", errReader(0))
			writer := httptest.NewRecorder()

			// Act
			handler.CreateTask(writer, req)

			// Assert
			assert.Equal(context, 500, writer.Code)
			assert.Equal(context, "error\n", writer.Body.String())
		})

		It("should return error if request body does not contain a task", func() {
			// Arrange
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader([]byte("")))
			writer := httptest.NewRecorder()

			// Act
			handler.CreateTask(writer, req)

			// Assert
			assert.Equal(context, 500, writer.Code)
			assert.NotNil(context, writer.Body.String())
		})

		It("should return error if storing task fails", func() {
			// Arrange
			directRedis.Close()
			taskString := `{"name": "test", "image": "alpine", "init": "init.sh"}`
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader([]byte(taskString)))
			writer := httptest.NewRecorder()

			// Act
			handler.CreateTask(writer, req)

			// Assert
			assert.Equal(context, 500, writer.Code)
			assert.NotNil(context, writer.Body.String())
		})

		It("should return task id if task creation succeeds", func() {
			// Arrange
			taskString := `{"name": "test", "image": "alpine", "init": "init.sh"}`
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader([]byte(taskString)))
			writer := httptest.NewRecorder()

			// Act
			handler.CreateTask(writer, req)

			// Assert
			assert.Equal(context, 201, writer.Code)
			id, err := uuid.FromString(writer.Body.String())
			assert.Nil(context, err)
			assert.NotNil(context, id)
		})

		It("should return error if adding to task queue fails", func() {
			// Arrange
			givenID := uuid.Must(uuid.NewV4())
			mockTaskStore := &mocks.Store{}
			mockTaskStore.On("StoreTask", mock2.AnythingOfType("model.Spec")).Return(&givenID, nil)
			mockTaskStore.On("Schedule", &givenID).Return(nil, errors.New("error"))

			handler = route.NewTaskHandlerImpl(mockTaskStore)
			taskString := `{"name": "test", "image": "alpine", "init": "init.sh"}`
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader([]byte(taskString)))
			writer := httptest.NewRecorder()

			// Act
			handler.CreateTask(writer, req)

			// Assert
			assert.Equal(context, 500, writer.Code)
			assert.Equal(context, "error\n", writer.Body.String())
		})
	})
})

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}
