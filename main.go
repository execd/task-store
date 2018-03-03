package taskman

import (
	"github.com/gorilla/mux"
	"github.com/wayofthepie/jobby-taskman/pkg/route"
)


func main() {
	router := route.ServiceImpl{Router: mux.NewRouter()}
	router.Listen()
}
