package main

import (
	"fmt"
	"net/http"
)

func logRequests(prefix string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(prefix, "----- REQUEST -----")
		fmt.Println(r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}
