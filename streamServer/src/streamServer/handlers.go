package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	VIDEO_DIR = "../videos/"
	MAX_UPLOAD_SIZE = 1024*1024*1024
)

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//这里我们可以通过拿到数据序列化成二进制，然后返回给web，比较复杂，这里用了比较简单的一种
	vid := p.ByName("vid-id")
	videoLink := VIDEO_DIR + vid

	video, err := os.Open(videoLink)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "internal Error")
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()

	targetUrl := ""
	http.Redirect(w,r,targetUrl,301)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//http.MaxBytesReader()限定了读取到缓冲区的值
	r.Body = http.MaxBytesReader(w, r.Body, int64(MAX_UPLOAD_SIZE))
	//检查大小，这里
	if err := r.ParseMultipartForm(int64(MAX_UPLOAD_SIZE)); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}

	// TODO
	// FormFile returns the first file for the provided form key.
	// FormFile calls ParseMultipartForm and ParseForm if necessary.
	// func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error){}
	// 可通过fileHeader进一步校验，也可以通过前端的 accept="video/*"来校验Content-Type
	file, _, err := r.FormFile("file") // <form name="file">
	if err != nil {
		log.Printf("Error when try to get file: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+fn, data, 0666) //no 777
	if err != nil {
		log.Printf("Write file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	/*ossfn := "videos/" + fn
	path := "./videos/" + fn
	bn := "aaron-video"
	ret := UploadToOss(ossfn, path, bn)
	if !ret {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	os.Remove(path)
	*/
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully")
}

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles(VIDEO_DIR + "upload.html")
	t.Execute(w, nil)

}


/*
1.  local -> oss   简单   消耗服务器的流量和带宽
2.	client -> local -> policy, client+policy ->oss  2和3都要开启跨域访问
3.	client -> local -> policy, client+policy ->oss ->callBack ->local sever
*/