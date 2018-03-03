package main

import (
	"github.com/gorilla/mux"
	"github.com/wayofthepie/task-store/pkg/route"
)


func main() {
	router := route.ServiceImpl{Router: mux.NewRouter()}
	router.Listen()
}
