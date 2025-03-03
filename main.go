package main

import (
	"authentification_service/server"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "2fa api")
	})
	api.HandleFunc("/check", server.CheckTOTPCodeHandler).Methods(http.MethodPost)
	api.HandleFunc("/users/get", server.GetUserRandomCodeHandler).Methods(http.MethodPost)
	api.HandleFunc("/users/add", server.AddUserHandler).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
