package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// Service interface for routes
type Service interface{}

// ServiceImpl is an implementation of Service
type ServiceImpl struct {
	Router *mux.Router
}

// Listen on the port defined by the env var PORT
func (s *ServiceImpl) Listen() {
	s.routes()
	if port, exists := os.LookupEnv("PORT"); exists {
		fmt.Printf("Running %s", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), s.Router))
	} else {
		log.Fatal("PORT env var not set")
	}

}

func (s *ServiceImpl) routes() {

}
