package main

import (
	"authentification_service/server"
	"authentification_service/storage"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	redis_cfg := storage.Config{
		Addr:        "localhost:6379",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	redis_client, err := storage.NewClient(context.Background(), redis_cfg)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "2fa api")
	})
	api.HandleFunc("/check", server.CheckTOTPCodeHandler(redis_client)).Methods(http.MethodPost)
	api.HandleFunc("/users/get", server.GetUserRandomCodeHandler).Methods(http.MethodPost)
	api.HandleFunc("/users/add", server.AddUserHandler).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
