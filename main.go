package main

import (
	"log"
	"net/http"

	"com.example/handlers"
	"github.com/gorilla/mux"
)

func main() {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/auth/{guid}", handlers.AuthUserHandler).Methods("GET")
	muxRouter.HandleFunc("/auth/refresh", handlers.RefreshHandler).Methods("PUT")

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", muxRouter))

}
