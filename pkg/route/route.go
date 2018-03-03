package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wayofthepie/task-store/pkg/task"
	"log"
	"net/http"
	"os"
)

// Service interface for routes
type Service interface {
	Listen()
}

// ServiceImpl is an implementation of Service
type ServiceImpl struct {
	router      *mux.Router
	taskHandler task.Handler
}

// NewServiceImpl : constructs a ServiceImpl
func NewServiceImpl(router *mux.Router, taskHandler task.Handler) *ServiceImpl {
	return &ServiceImpl{router: router, taskHandler: taskHandler}
}

// Listen on the port defined by the env var PORT
func (s *ServiceImpl) Listen() {
	s.routes()
	if port, exists := os.LookupEnv("PORT"); exists {
		fmt.Printf("Running %s", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), s.router))
	} else {
		log.Fatal("PORT env var not set")
	}

}

func (s *ServiceImpl) routes() {
	s.router.HandleFunc("/task/", s.taskHandler.CreateTaskHandler).Methods(http.MethodPost)
}
