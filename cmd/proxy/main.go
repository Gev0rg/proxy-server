package main

import (
	// "crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/Gev0rg/proxy-server/config"
	"github.com/Gev0rg/proxy-server/proxy"
	"github.com/Gev0rg/proxy-server/storage"
)

func main() {
	conf := config.NewConfig();

	store := &storage.Storage{}
	store.Connect()

	p := proxy.Proxy{
		Store: store,
	}

	server := http.Server{
		Handler: &p,
		Addr:    ":8080",
	}

	if conf.HTTPS {
		fmt.Println("Start serving TLS")
		if err := server.ListenAndServeTLS("certs/ca.crt", "certs/ca.key"); err != nil {
			log.Fatalln(err)
		}
	} else {
		fmt.Println("Start serving HTTP")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}
}
