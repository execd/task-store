package main

import (
	"github.com/execd/task-store/pkg/model"
	"github.com/execd/task-store/pkg/route"
	"github.com/execd/task-store/pkg/store"
	"github.com/execd/task-store/pkg/task"
	"github.com/execd/task-store/pkg/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := initializeRouter()
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}

func initializeRouter() *mux.Router {
	redis := store.NewClient("localhost:6379")
	uuidGen := util.NewUUIDGenImpl()
	taskStore := task.NewStoreImpl(redis, uuidGen)

	taskHandler := route.NewTaskHandlerImpl(taskStore, &model.Config{})

	router := mux.NewRouter()
	router.HandleFunc("/tasks/", taskHandler.CreateTask).Methods(http.MethodPost)

	return router
}
