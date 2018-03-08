package main

import (
	"github.com/gorilla/mux"
	"github.com/wayofthepie/task-store/pkg/event"
	"github.com/wayofthepie/task-store/pkg/route"
	"github.com/wayofthepie/task-store/pkg/store"
	"github.com/wayofthepie/task-store/pkg/task"
	"github.com/wayofthepie/task-store/pkg/model"
	"log"
)

func main() {
	initializeRouter().Listen()
}

func initializeRouter() route.Service {

	rabbit, _ := event.NewRabbitServiceImpl("amqp://guest:guest@localhost:5672/")
	taskSpec := &model.Spec{Image: "alpine", Name: "test", Init: "init.sh"}
	rabbit.PublishWork(taskSpec)

	redis := store.NewClient("localhost:6379")
	taskQueue := task.NewQueueImpl(redis)

	taskHandler := task.NewHandlerImpl(taskQueue, rabbit)
	router := route.NewServiceImpl(mux.NewRouter(), taskHandler)
	listener, _ := event.NewServiceImpl(rabbit)

	err := listener.ListenForTaskStatus()
	if err != nil {
		log.Fatalf("Failed when trying to listen for task statuses : %s", err.Error())
	}

	return router
}
