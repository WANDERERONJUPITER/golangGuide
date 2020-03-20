package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {

	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil  {
		p := &HomePage{
			Name:"aaronAnderson",
		}

		//todo nginx
		t, err :=template.ParseFiles("./templates/home.html")
		if err!=nil {
			log.Printf("parsing templates home.html failed: ", err)
		}
		t.Execute(w,p)
		return
	}

	//实际的操作中还要判断其是否匹配，交给前端通过jQuery等
	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userHome", http.StatusFound)
		return
	}

}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	//这里我们通过前端直接校验，确保传过来的数据是相对更正确的，减少服务端的压力
	fname := r.FormValue("username")

	var p *UserPage
	//先从cookie中读取，没有的话再从表单中读取
	if len(cname.Value) != 0 {
		p = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		p = &UserPage{Name: fname}
	}

	t, e := template.ParseFiles("./templates/userHome.html")
	if e != nil {
		log.Printf("Parsing userhome.html error: %s", e)
		return
	}

	t.Execute(w, p)
}

func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	//request 之前的预处理
	if r.Method != http.MethodPost {
		re, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(re))
		return
	}


	res, _ := ioutil.ReadAll(r.Body)
	apiBody := &ApiBody{}
	if err := json.Unmarshal(res, apiBody); err != nil {
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(re))
		return
	}

	//处理request
	request(apiBody, w, r)
	defer r.Body.Close()
}

// 域名转换，这里的整个header并没有改变
func proxyUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	//TODO 这里在实际中不应该写死
	u, _ := url.Parse("http://127.0.0.1:9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w,r)
}

//在前端中（home.scripts line:134  window.location.hostname），一旦我们将streamServer部署到其他服务器，将会出现跨域
//同proxyUploadHandler一样，我们考虑后期的扩展性，应该将不同的请求用不同的handler来处理
func proxyVideoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	u, _ := url.Parse("http://" + config.GetLbAddr() + ":9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}