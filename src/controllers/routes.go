package controllers

import (
	"GoMiddleMan/src/controllers/user"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", pong).Methods(http.MethodGet)
	userController.RegisterRoutes(router)
	return router
}

func pong(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	_,_ = fmt.Fprint(writer, "pong")
}