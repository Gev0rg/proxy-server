package main

import (
	// "crypto/tls"
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

	// fmt.Println("Start serving HTTPS")
	// if err := server.ListenAndServeTLS("certs/ca.crt", "certs/ca.key"); err != nil {
	// 	log.Fatalln(err);
	// }
}
