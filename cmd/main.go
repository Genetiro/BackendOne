package main

import (
	"fmt"
	"net/http"
	"time"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "GET")
	case http.MethodPost:
		fmt.Fprintln(w, "POST")
	}
}
func main() {
	handler := &Handler{}
	http.Handle("/", handler)
	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	srv.ListenAndServe()
}
