package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro/v2"
	proto "golangGuide/goMicro/src/helloworld/proto"
)

func main() {
	// Create a new service
	service := micro.NewService(micro.Name("label your service"))
	//service := micro.NewService()

	// Initialise the client and parse command line flags
	service.Init()

	// Create new greeter client
	greeter := proto.NewGreeterService("HelloWorld", service.Client())

	// Call the greeter
	rsp, err := greeter.Hello(context.TODO(), &proto.Request{Name: "aaron"})
	if err != nil {
		fmt.Println(err)
	}

	// Print response
	fmt.Println(rsp.Greeting)
}