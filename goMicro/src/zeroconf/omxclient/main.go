package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/grandcat/zeroconf"
)

func serviceCall(ip string, port int) {
	url := fmt.Sprintf("http://%v:%v", ip, port)

	log.Println("Making call to", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Got response: %s\n", data)
}

func main() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Channel to receive discovered service entries
	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println("Found service:", entry.ServiceInstanceName(), entry.Text)
			serviceCall(entry.AddrIPv4[0].String(), entry.Port)
		}
	}(entries)

	ctx := context.Background()

	err = resolver.Browse(ctx, "_omxremote._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
}