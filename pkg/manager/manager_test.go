package manager_test

import (
	"github.com/execd/task-store/mocks"
	"github.com/execd/task-store/pkg/manager"
	"github.com/execd/task-store/pkg/model"
	. "github.com/onsi/ginkgo"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var context = GinkgoT()

var _ = Describe("manage tasks", func() {
	Describe("manage task creation", func() {
		var taskStoreMock *mocks.Store
		var eventManagerMock *mocks.EventManager
		var taskManager *manager.TaskManagerImpl
		var quit chan int

		BeforeEach(func() {
			taskStoreMock = &mocks.Store{}
			eventManagerMock = &mocks.EventManager{}
			config := &model.Config{
				Manager: model.ManagerInfo{
					ExecutionQueueSize: 2,
					TaskQueueSize:      1,
				},
			}
			eventManagerMock.On("ListenForProgress", mock.Anything).Return(nil, nil)
			taskManager = manager.NewTaskManagerImpl(taskStoreMock, eventManagerMock, config)
			quit = make(chan int)
		})

		It("should not continue if there is a failure to check the execution set size", func() {
			// Arrange
			defer close(quit)
			taskStoreMock.On("ListenForTaskCreatedEvents").Return(buildCreatedTasksCh())
			taskStoreMock.On("ExecutingSetSize").Return(int64(0), errors.New("error"))

			// Act
			taskManager.ManageTasks(quit)

			// Assert
			taskStoreMock.AssertNotCalled(context, "GetTask", mock.Anything)
			quit <- 1
		})

		It("should not continue if execution queue is full", func() {
			// Arrange
			defer close(quit)
			taskStoreMock.On("ListenForTaskCreatedEvents").Return(buildCreatedTasksCh())
			taskStoreMock.On("ExecutingSetSize").Return(int64(2), nil)
			taskStoreMock.On("GetTask", mock.AnythingOfType("UUID"))

			// Act
			taskManager.ManageTasks(quit)

			// Assert
			taskStoreMock.AssertNotCalled(context, "AddTaskToExecutingSet")
			quit <- 1
		})

		It("should not continue if get task fails", func() {
			// Arrange
			defer close(quit)
			taskStoreMock.On("ListenForTaskCreatedEvents").Return(buildCreatedTasksCh())
			taskStoreMock.On("ExecutingSetSize").Return(int64(1), nil)
			taskStoreMock.On("GetTask", mock.AnythingOfType("*uuid.UUID")).Return(nil, errors.New("error"))

			// Act
			taskManager.ManageTasks(quit)

			// Assert
			eventManagerMock.AssertNotCalled(context, "PublishWork", mock.Anything)
			quit <- 1
		})

		It("should not continue if publish work fails", func() {
			// Arrange
			defer close(quit)
			taskStoreMock.On("ListenForTaskCreatedEvents").Return(buildCreatedTasksCh())
			taskStoreMock.On("ExecutingSetSize").Return(int64(1), nil)
			taskStoreMock.On("GetTask", mock.AnythingOfType("*uuid.UUID")).Return(nil, errors.New("error"))
			eventManagerMock.On("PublishWork", mock.Anything).Return(nil, errors.New("error"))

			// Act
			taskManager.ManageTasks(quit)

			// Assert
			taskStoreMock.AssertNotCalled(context, "AddTaskToExecutingSet", mock.Anything)
			quit <- 1
		})

		It("should not continue if adding to executing set fails", func() {
			// Arrange
			defer close(quit)
			taskStoreMock.On("ListenForTaskCreatedEvents").Return(buildCreatedTasksCh())
			taskStoreMock.On("ExecutingSetSize").Return(int64(1), nil)
			taskStoreMock.On("GetTask", mock.AnythingOfType("*uuid.UUID")).Return(nil, errors.New("error"))
			eventManagerMock.On("PublishWork", mock.Anything).Return(nil, errors.New("error"))
			taskStoreMock.On("AddTaskToExecutingSet",
				mock.AnythingOfType("*uuid.UUID")).Return(nil, errors.New("error"))

			// Act
			taskManager.ManageTasks(quit)
			quit <- 1
		})
	})
})

func buildCreatedTasksCh() <-chan *uuid.UUID {
	givenID := uuid.Must(uuid.NewV4())
	createdTasks := make(chan *uuid.UUID, 1)
	createdTasks <- &givenID
	return createdTasks
}
