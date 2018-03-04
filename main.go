package main

import (
	"github.com/gorilla/mux"
	"github.com/wayofthepie/task-store/pkg/event"
	"github.com/wayofthepie/task-store/pkg/model"
	"github.com/wayofthepie/task-store/pkg/route"
	"github.com/wayofthepie/task-store/pkg/store"
	"github.com/wayofthepie/task-store/pkg/task"
)

func main() {
	initializeRouter().Listen()
}

func initializeRouter() route.Service {
	conn, _ := event.NewRabbitConnection("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	amqp, _ := task.NewAmqpEventService(ch)
	taskSpec := &model.TaskSpec{Image: "alpine", Name: "test", Init: "init.sh"}
	amqp.PublishWork(taskSpec)
	redis := store.NewClient("localhost:6379")
	taskQueue := task.NewQueueImpl(redis)
	taskHandler := task.NewHandlerImpl(taskQueue)
	router := route.NewServiceImpl(mux.NewRouter(), taskHandler)
	return router
}
