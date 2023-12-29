package fileserver

import (
	"log"
	"net/http"
)

func Run(ipaddr string, filePath string) {

	serveVideo := makeVideoServer(filePath)

	http.HandleFunc("/", serveVideo)

	log.Printf("Starting serving of: \"%s\" @ %s", filePath, ipaddr)
	err := http.ListenAndServe(ipaddr+":80", nil)
	if err != nil {
		panic(err)
	}
}

func makeVideoServer(filePath string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Serving \"%s\" to %s", filePath, req.Host)
		http.ServeFile(w, req, filePath)
	}
}
