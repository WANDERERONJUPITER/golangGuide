package main

import (
	"github.com/julienschmidt/httprouter"
	"golangGuide/streamServer/src/api/sessions"
	"net/http"
)

//middleware
type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := &middleWareHandler{}
	m.r = r
	return m
}

// 如果没有这个方法（没有实现ServeHTTP），上面的NewMiddleWareHandler会报错，这里使用的就是duck type
func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check session
	m.r.ServeHTTP(w, r)
}

/*
	handler-->validation{1.request, 2.user} -->business logic -->response
	1.data model
	2.error handling
*/

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)
	//这里都会新开goroutine来处理，每个只占4k
	router.POST("/user/:username", Login)

	//这里路径写错了 导致调试一下午！！！
	router.GET("/user/:username", GetUserInfo)

	router.POST("/user/:username/videos", AddNewVideo)

	router.GET("/user/:username/videos", ListAllVideos)

	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)

	router.POST("/videos/:vid-id/comments", PostComment)

	router.GET("/videos/:vid-id/comments", ShowComments)
	return router
}

func Prepare() {
	sessions.LoadSessionsFromDB()
}

func main() {
	Prepare()
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mh)
}

//执行go install  其会生成在$PATH/bin目录下

/*
在router之前做一些诸如流控，鉴权，校验等处理，这一阶段通常称为middleware
main-->middleware-->defs(message,err)-->handlers-->DBops-->response
*/

//TODO   duck type
