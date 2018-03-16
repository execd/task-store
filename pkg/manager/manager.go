package manager

import (
	"fmt"
	"github.com/execd/task-store/pkg/model"
	"github.com/execd/task-store/pkg/task"
	"github.com/satori/go.uuid"
)

// TaskManager : manage tasks
type TaskManager interface {
	ManageTasks()
}

// TaskManagerImpl : a task manager impl
type TaskManagerImpl struct {
	store        task.Store
	eventManager task.EventManager
	config       *model.Config
}

// NewTaskManagerImpl : create a new task manager impl
func NewTaskManagerImpl(store task.Store, eventManager task.EventManager, config *model.Config) *TaskManagerImpl {
	return &TaskManagerImpl{store, eventManager, config}
}

// ManageTasks : manage task creation and progress
func (t *TaskManagerImpl) ManageTasks(quit <-chan int) {
	progressChQuit := make(chan int, 1)
	infoCh, errCh := t.eventManager.ListenForProgress(progressChQuit)
	go func() {
		for {
			select {
			case taskID := <-t.store.ListenForTaskCreatedEvents():
				t.scheduleForExecution(taskID)
			case info := <-infoCh:
				t.handleTaskProgressInfo(&info)
			case err := <-errCh:
				fmt.Printf("Received error from task progress watcher: %s\n", err.Error())
			case <-quit:
				progressChQuit <- 1
				return
			}
		}
	}()
}

func (t *TaskManagerImpl) scheduleForExecution(taskID *uuid.UUID) {
	size, err := t.store.ExecutingSetSize()
	if err != nil {
		fmt.Printf("Received error retrieving executing set size : %s", err.Error())
		return
	}

	if size >= t.config.Manager.ExecutionQueueSize {
		fmt.Println("Not scheduling taskSpec taskID for execution, executing set has reached capacity.")
		return
	}

	taskSpec, err := t.store.GetTask(taskID)
	if err != nil {
		fmt.Printf("Failed scheduling taskSpec taskID for execution: %s\n", err.Error())
		return
	}

	err = t.eventManager.PublishWork(taskSpec)
	if err != nil {
		fmt.Printf("Failed scheduling taskSpec taskID for execution: %s\n", err.Error())
		return
	}

	err = t.store.AddTaskToExecutingSet(taskID)
	if err != nil {
		fmt.Printf("Task spec scheduled for execution, but failed to add to executing set: %s\n", err.Error())
		return
	}

	fmt.Printf("Task %s successfully added to executing set\n", taskID.String())
}

func (t *TaskManagerImpl) handleTaskProgressInfo(info *model.Info) {
	err := t.store.UpdateTaskInfo(info)
	if err != nil {
		fmt.Printf("Received error trying to update task info: %s\n", err.Error())
		return
	}
	err = t.store.RemoveTaskFromExecutingSet(info.ID)
	if err != nil {
		fmt.Printf("Error removing task from executing set: %s\n", err.Error())
		return
	}
}
