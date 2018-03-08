package task

import (
	"github.com/alicebob/miniredis"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/wayofthepie/task-store/pkg/model"
	"github.com/wayofthepie/task-store/pkg/store"
	"fmt"
	"log"
)

var context = GinkgoT()

var _ = Describe("QueueImpl", func() {
	var taskQueue *QueueImpl
	var directRedis *miniredis.Miniredis

	BeforeEach(func() {
		s, err := miniredis.Run()
		if err != nil {
			panic(err)
		}

		redis := store.NewClient(s.Addr())
		taskQueue = &QueueImpl{redis: redis}
		directRedis = s
	})

	AfterEach(func() {
		defer directRedis.Close()
	})

	Describe("Push", func() {

		expectedTaskSpec := &model.Spec{Image: "alpine", Name: "test", Init: "init.sh"}

		It("should store task details as separate key", func() {
			// Act
			id, err := taskQueue.Push(expectedTaskSpec)
			failOnError(err)
			expectedTaskSpec.ID = id

			// Assert
			taskData, _ := directRedis.Get(fmt.Sprintf("task:%s", id.String()))
			taskSpec := new(model.Spec)
			taskSpec.UnmarshalBinary([]byte(taskData))

			assert.Nil(context, err)
			assert.Equal(context, expectedTaskSpec, taskSpec)
		})
	})
})

func failOnError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
