package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type KubernetesController struct {
	router *mux.Router
}

func (k *KubernetesController) Add() {
	log.Println("adding kuberneters controller")
	k.router.HandleFunc("/alive", k.getAlive)
	k.router.HandleFunc("/ready", k.getReady)
}

func (a *KubernetesController) getAlive(w http.ResponseWriter, _ *http.Request) {
	log.Println("kuberneters controller alive")
}

func (a *KubernetesController) getReady(w http.ResponseWriter, _ *http.Request) {
	log.Println("kuberneters controller ready")
}

func NewKubernetesController(router *mux.Router) *KubernetesController {
	return &KubernetesController{router: router}
}
