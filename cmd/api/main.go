package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Gev0rg/proxy-server/api"
	"github.com/Gev0rg/proxy-server/storage"
)
func main() {
	defer func() {
        if r := recover(); r != nil {
            // Обработка ошибки паники
            fmt.Println("Возникла паника:", r)
        }
    }()
	
	store := &storage.Storage{}
	store.Connect()

	handlers := &api.Handlers{
		Storage: store,
	}

	router := mux.NewRouter()
	router.HandleFunc("/requests", handlers.GetRequests).Methods(http.MethodGet)
	router.HandleFunc("/requests/{id}", handlers.GetRequestByID).Methods(http.MethodGet)
	router.HandleFunc("/repeat/{id}", handlers.RepeatRequest).Methods(http.MethodGet)
	router.HandleFunc("/dirsearch/{id}", handlers.DirSearch).Methods(http.MethodGet)

	server := http.Server{
		Handler: router,
		Addr:    ":8000",
	}

	fmt.Println("Start serving api")
	server.ListenAndServe()
}