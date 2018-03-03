package main

import (
	"github.com/gorilla/mux"
	"github.com/wayofthepie/task-store/pkg/route"
	"github.com/wayofthepie/task-store/pkg/store"
	"github.com/wayofthepie/task-store/pkg/task"
)

func main() {
	initializeRouter().Listen()
}

func initializeRouter() route.Service {
	redis := store.NewClient("localhost:6379")
	taskQueue := task.NewQueueImpl(redis)
	taskHandler := task.NewHandlerImpl(taskQueue)
	router := route.NewServiceImpl(mux.NewRouter(), taskHandler)
	return router
}
