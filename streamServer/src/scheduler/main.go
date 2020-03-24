package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video-delete-record/:vid-id", vidDelRecHandler)

	return router
}

func main() {
	//c :=make(chan int)
	go taskrunner.Start()
	r := RegisterHandlers()
	//<-c
	//这里会阻塞  上面可以放心开goroutine ,或者通过channel构造block区间
	http.ListenAndServe(":9001", r)
}
