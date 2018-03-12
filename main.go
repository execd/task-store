package main

import (
	"github.com/gorilla/mux"
	"github.com/wayofthepie/task-store/pkg/route"
	"github.com/wayofthepie/task-store/pkg/store"
	"github.com/wayofthepie/task-store/pkg/task"
	"github.com/wayofthepie/task-store/pkg/util"
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

	taskHandler := route.NewTaskHandlerImpl(taskStore)

	router := mux.NewRouter()
	router.HandleFunc("/tasks/", taskHandler.CreateTask).Methods(http.MethodPost)

	return router
}
