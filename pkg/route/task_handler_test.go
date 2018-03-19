package route_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alicebob/miniredis"
	"github.com/execd/task-store/mocks"
	"github.com/execd/task-store/pkg/model"
	"github.com/execd/task-store/pkg/redis"
	"github.com/execd/task-store/pkg/route"
	"github.com/execd/task-store/pkg/task"
	"github.com/execd/task-store/pkg/util"
	. "github.com/onsi/ginkgo"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
)

var context = GinkgoT()

var _ = Describe("task handler", func() {
	var taskStore *task.StoreImpl
	var directRedis *miniredis.Miniredis
	var handler *route.TaskHandlerImpl

	BeforeEach(func() {
		s, err := miniredis.Run()
		if err != nil {
			panic(err)
		}
		directRedis = s
		redis := redis.NewClient(s.Addr())
		uuidGen := util.NewUUIDGenImpl()
		taskStore = task.NewStoreImpl(redis, uuidGen)
		config := &model.Config{
			Manager: model.ManagerInfo{
				ExecutionQueueSize: 10,
				TaskQueueSize:      10,
			},
		}
		handler = route.NewTaskHandlerImpl(taskStore, config)
	})

	Describe("create task", func() {
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

		It("should return error of task queue size is greater than max task queue size", func() {
			// Arrange
			config := &model.Config{
				Manager: model.ManagerInfo{
					ExecutionQueueSize: 10,
					TaskQueueSize:      0,
				},
			}
			handler = route.NewTaskHandlerImpl(taskStore, config)
			givenID := uuid.Must(uuid.NewV4())
			taskStore.PushTask(&givenID)
			taskString := `{"name": "test", "image": "alpine", "init": "init.sh"}`
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader([]byte(taskString)))
			writer := httptest.NewRecorder()

			// Act
			handler.CreateTask(writer, req)

			// Assert
			assert.Equal(context, 500, writer.Code)
			assert.Equal(context, writer.Body.String(), "Failed to create task, task queue has reached its limit!\n")
		})
	})

	Describe("get task", func() {
		It("should return the task related to the given id", func() {
			// Arrange
			givenTaskSpec := model.Spec{
				Image:    "alpine",
				Init:     "init.sh",
				InitArgs: []string{"10"},
			}
			id, err := taskStore.StoreTask(givenTaskSpec)
			assert.Nil(context, err)
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader(nil))

			writer := httptest.NewRecorder()
			givenTaskSpec.ID = id
			vars := map[string]string{
				"id": id.String(),
			}
			// Act
			handler.GetTask(writer, req, vars)

			// Assert
			taskSpec := new(model.Spec)
			json.Unmarshal(writer.Body.Bytes(), taskSpec)
			assert.Equal(context, 200, writer.Code)
			assert.Equal(context, givenTaskSpec, *taskSpec)
		})

		It("should return error if id is not a v4 uuid", func() {
			// Arrange
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader(nil))
			writer := httptest.NewRecorder()

			vars := map[string]string{
				"id": "1234",
			}
			// Act
			handler.GetTask(writer, req, vars)

			// Assert
			assert.Equal(context, 500, writer.Code)
			assert.Contains(context, writer.Body.String(), "failed to build id from 1234")
		})

		It("should return error if task retrieval fails", func() {
			// Arrange
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader(nil))
			writer := httptest.NewRecorder()

			vars := map[string]string{
				"id": uuid.Must(uuid.NewV4()).String(),
			}
			directRedis.Close()

			// Act
			handler.GetTask(writer, req, vars)

			// Assert
			assert.Equal(context, 500, writer.Code)
			assert.NotNil(context, writer.Body.String())
		})

		It("should return error retrieved if task fails to marshal", func() {
			// Arrange
			req, _ := http.NewRequest("POST", "/handle", bytes.NewReader(nil))
			writer := httptest.NewRecorder()
			givenID := uuid.Must(uuid.NewV4())
			vars := map[string]string{
				"id": givenID.String(),
			}
			taskStoreMock := &mocks.Store{}
			config := &model.Config{
				Manager: model.ManagerInfo{
					ExecutionQueueSize: 10,
					TaskQueueSize:      10,
				},
			}
			handler = route.NewTaskHandlerImpl(taskStoreMock, config)

			taskStoreMock.On("GetTask", mock.Anything).Return(nil, errors.New("error"))
			// Act
			handler.GetTask(writer, req, vars)

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
