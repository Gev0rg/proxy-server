package main

import (
	// "crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Gev0rg/proxy-server/api"
	"github.com/Gev0rg/proxy-server/proxy"
	"github.com/Gev0rg/proxy-server/storage"
)

func main() {
	store := &storage.Storage{}
	store.Connect()

	p := proxy.Proxy{
		Store: store,
	}

	handlers := &api.Handlers{
		Storage: store,
		Proxy:   p,
	}

	server := http.Server{
		Handler: &p,
		Addr:    ":8080",
	}

	router := mux.NewRouter()
	router.HandleFunc("/requests", handlers.GetRequests).Methods(http.MethodGet)
	router.HandleFunc("/requests/{id}", handlers.GetRequestByID).Methods(http.MethodGet)
	router.HandleFunc("/repeat/{id}", handlers.RepeatRequest).Methods(http.MethodGet)
	router.HandleFunc("/dirsearch/{id}", handlers.DirSearch).Methods(http.MethodGet)

	apiServer := http.Server{
		Handler: router,
		Addr:    ":8000",
	}

	// fmt.Println("Start serving HTTP")
	// if err := server.ListenAndServe(); err != nil {
	// 	log.Fatalln(err)
	// }

	fmt.Println("Start serving HTTPS")
	go apiServer.ListenAndServe()
	if err := server.ListenAndServeTLS("certs/ca.crt", "certs/ca.key"); err != nil {
		log.Fatalln(err)
	}
}
