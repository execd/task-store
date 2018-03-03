package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"os"
	"log"
)

type Service interface {
}

type ServiceImpl struct {
	Router *mux.Router
}

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
