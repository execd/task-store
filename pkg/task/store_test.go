package task_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/alicebob/miniredis"
	"github.com/wayofthepie/task-store/pkg/store"
	"github.com/satori/go.uuid"
	"github.com/wayofthepie/task-store/pkg/model"
	"github.com/wayofthepie/task-store/pkg/task"
	"github.com/stretchr/testify/assert"
	"fmt"
)

var _ = Describe("task store", func() {
	var taskStore *task.StoreImpl
	var directRedis *miniredis.Miniredis
	var givenID uuid.UUID
	var givenTaskSpec model.Spec

	BeforeEach(func() {
		s, err := miniredis.Run()
		if err != nil {
			panic(err)
		}

		redis := store.NewClient(s.Addr())
		taskStore = task.NewStoreImpl(redis)
		directRedis = s
		givenID = uuid.Must(uuid.NewV4())
		givenTaskSpec = model.Spec{
			ID:       &givenID,
			Name:     "test",
			Image:    "alpine",
			Init:     "init.sh",
			InitArgs: []string{"10"},
		}
	})

	AfterEach(func() {
		defer directRedis.Close()
	})

	Describe("storing a task", func() {
		It("should return task id", func() {
			// Act
			id, err := taskStore.CreateTask(givenTaskSpec)

			// Assert
			assert.Nil(context, err)
			assert.Equal(context, &givenID, id)
		})

		It("should store the given task by its id", func() {
			// Arrange
			assert.Empty(context, directRedis.Keys())

			// Act
			_, err := taskStore.CreateTask(givenTaskSpec)

			// Assert
			assert.Nil(context, err)
			assert.Len(context, directRedis.Keys(), 1)
			assert.Contains(context, directRedis.Keys()[0], givenID.String())
		})

		It("should return error if task with id already exists", func() {
			// Arrange
			_, err := taskStore.CreateTask(givenTaskSpec)
			assert.Nil(context, err)

			// Act
			_, err = taskStore.CreateTask(givenTaskSpec)

			// Assert
			assert.NotNil(context, err)
			assert.Equal(context, fmt.Sprintf("task with id %s already exists", givenID), err.Error())
		})

		It("should return error if storing task fails", func() {
			// Arrange
			directRedis.Close()

			// Act
			_, err := taskStore.CreateTask(givenTaskSpec)

			// Assert
			assert.NotNil(context, err)
			assert.Equal(context, fmt.Sprintf("storing task with id %s failed", givenID), err.Error())
		})
	})

	Describe("retrieving  a task", func() {
		It("should return the task related to the given id", func () {
			// Arrange
			_, err := taskStore.CreateTask(givenTaskSpec)
			failOnError(err)

			// Act
			task, err := taskStore.GetTask(givenID)

			// Assert
			assert.Nil(context, err)
			assert.Equal(context, &givenTaskSpec, task)
		})

		It("should return an error if task with given id does not exist", func () {
			// Arrange
			assert.Empty(context, directRedis.Keys())

			// Act
			_, err := taskStore.GetTask(givenID)

			// Assert
			assert.NotNil(context, err)
			assert.Equal(context, fmt.Sprintf("failed to retrieve task with id %s", givenID.String()), err.Error())
		})

		It("should return an error if retrieved task fails to be re-constructed", func () {
			// Arrange
			badData := "test"
			_, err := taskStore.CreateTask(givenTaskSpec)
			failOnError(err)
			taskID := directRedis.Keys()[0]
			directRedis.Set(taskID, badData)

			// Act
			_, err = taskStore.GetTask(givenID)

			// Assert
			assert.NotNil(context, err)
			assert.Equal(context,
				fmt.Sprintf("failed to build task with id %s from retrieved data %s", givenID.String(), badData),
					err.Error())
		})
	})
})
