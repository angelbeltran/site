package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	staticFileServer := logRequests("static", http.FileServer(http.Dir("./static")))
	mux.Handle("/css/", staticFileServer)
	mux.Handle("/js/", staticFileServer)
	mux.Handle("/images/", staticFileServer)
	mux.Handle("/favicon.ico", staticFileServer)
	mux.Handle("/", logRequests("html", &htmlTemplateServer{
		root: "./static/html",
	}))

	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal("server error:", err)
	}
}
