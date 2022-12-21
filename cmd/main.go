package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Genetiro/BackendOne/tree/project/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	handler := router.HandleFunc("/", (Home))

	srv := server.New(":8080", handler)

	srv.Start()

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-osSigChan
	log.Println("received OS interrupting signal")
	srv.Stop()
	router.Post("./internal/html/home.tmpl", http.HandlerFunc(Home))
	srv.ListenAndServe()
}
