package main

import (
	"github.com/execd/task-store/pkg/config"
	"github.com/execd/task-store/pkg/event"
	"github.com/execd/task-store/pkg/manager"
	"github.com/execd/task-store/pkg/model"
	"github.com/execd/task-store/pkg/route"
	"github.com/execd/task-store/pkg/store"
	"github.com/execd/task-store/pkg/task"
	"github.com/execd/task-store/pkg/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		panic("You must give the config file location as an argument!")
	}

	configFile := args[0]
	conf := parseConfig(configFile)

	taskStore := initializeStore()

	initializeAndLaunchManager(taskStore, conf)

	router := initializeRouter(taskStore, conf)
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}

func initializeAndLaunchManager(taskStore task.Store, config *model.Config) {
	rabbit, err := event.NewRabbitServiceImpl("amqp://localhost:5672")
	if err != nil {
		panic(err.Error())
	}
	eventManager, err := task.NewEventManagerImpl(rabbit)
	if err != nil {
		panic(err.Error())
	}

	taskManager := manager.NewTaskManagerImpl(taskStore, eventManager, config)
	quit := make(chan int)
	taskManager.ManageTasks(quit)
}

func initializeStore() task.Store {
	redis := store.NewClient("localhost:6379")
	uuidGen := util.NewUUIDGenImpl()
	return task.NewStoreImpl(redis, uuidGen)
}

func initializeRouter(taskStore task.Store, config *model.Config) *mux.Router {
	taskHandler := route.NewTaskHandlerImpl(taskStore, config)
	router := mux.NewRouter()

	router.HandleFunc("/tasks/", taskHandler.CreateTask).Methods(http.MethodPost)
	getTaskH := func(w http.ResponseWriter, r *http.Request) {
		taskHandler.GetTask(w, r, mux.Vars(r))
	}
	router.HandleFunc("/tasks/{id}", getTaskH).Methods(http.MethodGet)

	return router
}

func parseConfig(path string) *model.Config {
	parser := config.NewParserImpl()
	return parser.ParseConfig(path)
}
