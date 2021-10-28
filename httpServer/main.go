package main

import (
	"flag"
	"net/http"
	"os"
	"regexp"

	"github.com/golang/glog"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`:.*$`)
	clientIp := string(re.ReplaceAll([]byte(r.RemoteAddr), []byte("")))
	defer func() {
		err := recover()
		if err != nil {
			glog.V(2).Infof("request=rootHandler, remoteAddr=%s, statusCode=%d", clientIp, http.StatusInternalServerError)
		} else {
			glog.V(2).Infof("request=rootHandler, remoteAddr=%s, statusCode=%d", clientIp, http.StatusOK)
		}

	}()

	writeHeader(w, r.Header)
	w.Write([]byte("Hello"))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("200"))
}

func writeHeader(w http.ResponseWriter, h http.Header) {
	for k, v := range h {
		for _, v := range v {
			w.Header().Set(k, v)
		}
	}
	version := os.Getenv("VERSION")
	if version != "" {
		w.Header().Add("Version", version)
	}
}

func main() {
	flag.Set("v", "4")
	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	glog.V(2).Info("Starting http server")

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)

	err := http.ListenAndServe(":8080", mux)
	glog.Error(err.Error())
}
