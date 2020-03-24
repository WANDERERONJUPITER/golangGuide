package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

/*
TODO  compare to next one， which one is better
func (mdh *middleWareHandler) New(r *httprouter.Router,cc int) http.Handler  {
	m :=  middleWareHandler{}
	m.r =r
	m.l = newConnLimiter(cc)
	return m
}*/

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = newConnLimiter(cc)
	return m
}
func (mdh middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !mdh.l.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "too many requests") //sc 429
		return
	}
	//TODO  这里第一次直接写成了这个 直接导致连接 too may requests
	//mdh.ServeHTTP(w,r)
	mdh.r.ServeHTTP(w, r)

	defer mdh.l.ReleaseConn()
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/videos/:vid-id", streamHandler)
	router.POST("/upload/:vid-id", uploadHandler)
	router.GET("/testPage", testPageHandler)
	return router
}

func main() {
	r := RegisterHandlers()
	mdh := NewMiddleWareHandler(r, 2)
	http.ListenAndServe(":9000", mdh)
}
