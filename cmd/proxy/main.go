package main

import (
	// "crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/Gev0rg/proxy-server/proxy"
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

	p := proxy.Proxy{
		Store: store,
	}

	server := http.Server{
		Handler: &p,
		Addr:    ":8080",
	}

	fmt.Println("Start serving HTTP")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

	// fmt.Println("Start serving HTTPS")
	// if err := server.ListenAndServeTLS("certs/ca.crt", "certs/ca.key"); err != nil {
	// 	log.Fatalln(err)
	// }
}
