package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/grandcat/zeroconf"
)

// Our fake service.
// This could be a HTTP/TCP service or whatever you want.
func startService() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Hello world!")
	})

	log.Println("starting http service...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Start out http service
	go startService()

	// Extra information about our service
	meta := []string{
		"version=0.1.0",
		"hello=world",
	}

	service, err := zeroconf.Register(
		"awesome-sauce",   // service instance name
		"_omxremote._tcp", // service type and protocl
		"local.",          // service domain
		8080,              // service port
		meta,              // service metadata
		nil,               // register on all network interfaces
	)

	if err != nil {
		log.Fatal(err)
	}

	defer service.Shutdown()

	// Sleep forever
	select{}
}