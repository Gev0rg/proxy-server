package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gev0rg/proxy-server/proxy"
)

func main() {
	p := proxy.Proxy{}

	server := http.Server{
		Handler: 	&p,
		Addr: 		":8080",
	}

	fmt.Println("Start serving HTTP")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err);
	}
}
