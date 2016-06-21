package webcam

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	proxy httputil.ReverseProxy
}

func NewHandler(
	webcamHost string,
) Handler {
	log.Printf("Webcam handler %v", webcamHost)

	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = webcamHost
		req.URL.Path = "/"
		req.URL.RawQuery = "action=stream"
	}

	flushInterval, err := time.ParseDuration("10ms")
	if err != nil {
		log.Fatal("golang broke", err)
	}

	proxy := httputil.ReverseProxy{
		Director:      director,
		FlushInterval: flushInterval,
		ErrorLog:      log.New(ioutil.Discard, "", 0),
	}

	return &handler{
		proxy: proxy,
	}
}

func (h handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}
