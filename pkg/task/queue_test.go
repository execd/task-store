package task_test

import (
	"github.com/alicebob/miniredis"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/wayofthepie/task-store/pkg/store"
	"log"
	"github.com/satori/go.uuid"
	"github.com/wayofthepie/task-store/pkg/task"
)

var context = GinkgoT()

var _ = Describe("queue", func() {
	var taskQueue *task.QueueImpl
	var directRedis *miniredis.Miniredis

	BeforeEach(func() {
		s, err := miniredis.Run()
		if err != nil {
			panic(err)
		}

		redis := store.NewClient(s.Addr())
		taskQueue = task.NewQueueImpl(redis)
		directRedis = s
	})

	AfterEach(func() {
		defer directRedis.Close()
	})

	Describe("pushing a task on the queue", func() {
		It("should add task to task queue", func() {
			// Arrange
			givenId := uuid.Must(uuid.NewV4())

			// Act
			_, err := taskQueue.Push(&givenId)
			failOnError(err)

			// Assert
			taskIDs, err := directRedis.List("taskQ")
			failOnError(err)

			assert.Nil(context, err)
			assert.Len(context, taskIDs, 1)
			assert.Contains(context, taskIDs, givenId.String())
		})

		It("should return pushed task id", func () {
			// Arrange
			givenId := uuid.Must(uuid.NewV4())

			// Act
			id, err := taskQueue.Push(&givenId)
			failOnError(err)

			// Assert
			assert.Nil(context, err)
			assert.Equal(context, givenId, *id)
		})

		It("should return error if pushing task onto queue fails", func () {
			//Arrange
			directRedis.Close()
			givenId := uuid.Must(uuid.NewV4())

			// Act
			_, err := taskQueue.Push(&givenId)

			// Assert
			assert.NotNil(context, err)
		})
	})
})

func failOnError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
