package main

import (
	"encoding/json"
	"golangGuide/streamServer/src/api/defs"

	"io"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrResponse) {
	w.WriteHeader(errResp.HttpSC)
	resStr, _ := json.Marshal(&errResp.Error)

	io.WriteString(w, string(resStr))
}

//区别于上一个，这个的sc 可能会随着业务发生改变
func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
