package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

/*
将用户的前端的行为转换为我们固定的api的请求，如点击事件等
转发
proxy和api模式   用js来获取ajax的请求都会出现跨域
cross origin resource sharing   当跨域请求出现时候我们要使用proxy
如 前端访问 http://127.0.0.1:8000/upload/:vid-id
我们可以在后端可以转为： http://127.0.0.1:9000/upload/:vid-id
*/
var httpClient *http.Client
var (
	LBaddr  = ""
	OSSaddr = ""
)

func init() {
	httpClient = &http.Client{}
}

func request(b *ApiBody, w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error

	u, _ := url.Parse(b.Url)
	u.Host = LBaddr + ":" + u.Port()
	newUrl := u.String()

	switch b.Method {
	case http.MethodGet:
		req, _ := http.NewRequest("GET", newUrl, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST", newUrl, bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("DELETE", newUrl, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad api request")
		return
	}
}

func normalResponse(w http.ResponseWriter, r *http.Response) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		re, _ := json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w, string(re))
		return
	}

	w.WriteHeader(r.StatusCode)
	io.WriteString(w, string(res))
}
