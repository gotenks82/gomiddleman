package main

import (
	"GoMiddleMan/src/controllers"
	"GoMiddleMan/src/middleman"
	"log"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: controllers.GetRouter(),
	}
	defer server.Close()
	log.Println("Listening...")

	mm := middleman.GetInstance()

	if mm != nil {
		defer mm.Shutdown()
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}
