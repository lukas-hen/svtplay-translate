package server

import (
	"log"
	"net/http"
)

func RunServer() {
	http.HandleFunc("/", ServeVideo)
	err := http.ListenAndServe("192.168.1.106:80", nil)
	if err != nil {
		panic(err)
	}
}

func ServeVideo(w http.ResponseWriter, req *http.Request) {
	log.Print("Serving video ./out_subs.mp4")
	http.ServeFile(w, req, "./out_subs.mp4")
}
