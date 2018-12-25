package userController

import (
	"GoMiddleMan/src/middleman"
	"GoMiddleMan/src/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RegisterRoutes(router *mux.Router) {
	log.Print("registering user routes")
	router.HandleFunc("/user/{userId}/interest", addInterest).Methods(http.MethodPost)
}

func addInterest(writer http.ResponseWriter, request *http.Request) {
	var interest models.Interest
	err := json.NewDecoder(request.Body).Decode(&interest)

	params := mux.Vars(request)
	user := params["userId"]

	if err != nil {
		log.Printf("Error: %s", err)
		http.Error(writer, "Can't parse request", http.StatusBadRequest)
	}
	log.Printf("Received: %s, for UserId: %s", interest, user)

	middleman.GetInstance().AddInterest(user, interest)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(201)
	_ = json.NewEncoder(writer).Encode(interest)
}
