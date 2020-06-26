package api

import (
	"log"
	"net/http"
	"redis-proxy/src/api/controllers"

	"github.com/gorilla/mux"
)

type Router struct {
	port *string
}

func (r *Router) controllers(router *mux.Router) {
	controllers.NewKubernetesController(router).Add()
	controllers.NewRedisController(router).Add()
}

func (r *Router) Start() {
	log.Println("starting router", *r.port)

	router := mux.NewRouter().StrictSlash(true)

	r.controllers(router)

	go http.ListenAndServe(*r.port, router)
}

func NewRouter(port *string) *Router {
	return &Router{port: port}
}
