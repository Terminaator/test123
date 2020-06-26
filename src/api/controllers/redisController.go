package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"redis-proxy/src/resp"
	"redis-proxy/src/resp/constants"

	"github.com/gorilla/mux"
)

type RedisController struct {
	router *mux.Router
}

type Response struct {
	Response string
}

func (r *RedisController) getBuilding(w http.ResponseWriter, _ *http.Request) {
	r.request(&w, &constants.BUILDING_CODE)
}

func (r *RedisController) getUtilityBuilding(w http.ResponseWriter, _ *http.Request) {
	r.request(&w, &constants.UTILITY_BUILDING_CODE)
}

func (r *RedisController) getProcedure(w http.ResponseWriter, _ *http.Request) {
	r.request(&w, &constants.PROCEDURE_CODE)
}

func (r *RedisController) getBuildingPart(w http.ResponseWriter, _ *http.Request) {
	r.request(&w, &constants.BUILDING_PART_CODE)
}

func (r *RedisController) getDocument(w http.ResponseWriter, req *http.Request) {
	doty := mux.Vars(req)["doty"]
	command := []byte(fmt.Sprintf(constants.DOCUMENT_CODE, len(doty), doty))
	r.request(&w, &command)
}

func (r *RedisController) response(w *http.ResponseWriter, buf *[]byte) {
	if (*buf)[0] == '-' {
		(*w).WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(*w).Encode(Response{string((*buf)[1 : len(*buf)-2])})

}

func (r *RedisController) Do(buf *[]byte) *[]byte {
	redis := resp.NewRedis()

	defer redis.Close()

	return redis.Do(buf)
}

func (r *RedisController) request(w *http.ResponseWriter, buf *[]byte) {
	r.response(w, r.Do(buf))
}

func (r *RedisController) Add() {
	log.Println("adding redis controller")
	r.router.HandleFunc("/building", r.getBuilding)
	r.router.HandleFunc("/utilitybuilding", r.getUtilityBuilding)
	r.router.HandleFunc("/procedure", r.getProcedure)
	r.router.HandleFunc("/buildingpart", r.getBuildingPart)
	r.router.HandleFunc("/document/{doty}", r.getDocument)
}

func NewRedisController(router *mux.Router) *RedisController {
	return &RedisController{router: router}
}
