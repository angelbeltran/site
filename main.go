package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	publicFileServer := logRequests("public", http.FileServer(http.Dir("./public")))
	mux.Handle("/css/", publicFileServer)
	mux.Handle("/js/", publicFileServer)
	mux.Handle("/images/", publicFileServer)
	mux.Handle("/favicon.ico", publicFileServer)
	mux.Handle("/", logRequests("html", &htmlTemplateServer{
		root: "./public/html",
	}))

	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal("server error:", err)
	}
}
