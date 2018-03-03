package task

import (
	. "github.com/onsi/ginkgo"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/wayofthepie/jobby-taskman/pkg/store"
	"github.com/wayofthepie/jobby-taskman/pkg/model"
)

var context = GinkgoT()

var _ = Describe("TaskQueue", func() {
	var taskQueue *TaskQueueImpl
	var directRedis *miniredis.Miniredis

	BeforeEach(func() {
		s, err := miniredis.Run()
		if err != nil {
			panic(err)
		}

		redis := store.NewClient(s.Addr())
		taskQueue = &TaskQueueImpl{redis: redis}
		directRedis = s

	})

	AfterEach(func() {
		defer directRedis.Close()
	})

	Describe("GetTaskInfo", func() {
		It("should retrieve a task spec when one exists with the given id", func() {
			// Arrange
			taskId := "task:1"
			expectedTaskSpec := &model.TaskSpec{Image: "alpine", Name: "test", Init: "init.sh"}
			taskQueue.redis.Set(taskId, expectedTaskSpec, 0)

			// Act
			taskSpec, err := taskQueue.GetTaskInfo(taskId)

			// Assert
			assert.Nil(context, err)
			assert.Equal(context, expectedTaskSpec, taskSpec)
		})
	})

	Describe("Push", func() {

		expectedTaskSpec := &model.TaskSpec{Image: "alpine", Name: "test", Init: "init.sh"}

		It("should return the taskQueue length after successful push", func() {
			// Act
			result, err := taskQueue.Push(expectedTaskSpec)

			// Assert
			assert.Nil(context, err)
			assert.Equal(context, "task:1", result)
		})

		It("should store task details as separate key", func() {
			// Act
			_, err := taskQueue.Push(expectedTaskSpec)

			// Assert
			taskData, _ := directRedis.Get("task:1")
			taskSpec := new(model.TaskSpec)
			taskSpec.UnmarshalBinary([]byte(taskData))

			assert.Nil(context, err)
			assert.Equal(context, expectedTaskSpec, taskSpec)
		})

		It("should store successive tasks with incrementing id's", func() {
			// Arrange
			_, err := taskQueue.Push(expectedTaskSpec)

			// Act
			result, err := taskQueue.Push(expectedTaskSpec)

			// Assert
			taskData, _ := directRedis.Get("task:2")

			taskSpec := new(model.TaskSpec)
			taskSpec.UnmarshalBinary([]byte(taskData))

			assert.Nil(context, err)
			assert.Equal(context, "task:2", result)
			assert.Equal(context, expectedTaskSpec, taskSpec)
		})

		It("should return error if next task id cannot be reserved", func() {
			// Arrange
			directRedis.Set("task:id", "test")

			// Act
			_, err := taskQueue.Push(expectedTaskSpec)

			// Assert
			assert.NotNil(context, err)
			assert.Contains(context, err.Error(), "reserving an id failed with")
		})
	})
})
