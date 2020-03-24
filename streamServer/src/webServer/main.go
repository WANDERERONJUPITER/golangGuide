package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", homeHandler)
	router.POST("/", homeHandler)

	router.GET("/userHome", userHomeHandler)
	router.POST("/userHome", userHomeHandler)

	router.POST("/api", apiHandler)

	router.GET("/videos/:vid-id", proxyVideoHandler)

	//这里我们将streamServer中的服务在这里通过proxy实现
	router.POST("/upload/:vid-id", proxyUploadHandler)

	//静态文件,注意这里的目录写法   /statics    不加/会出现严重的报错
	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))

	return router
}

func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}
