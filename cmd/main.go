package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	srv := NewSrv("8080", r)

	srv.Start()

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-osSigChan
	log.Println("received OS interrupting signal")
	srv.Stop()
	r.Mount("/links", LinkResource{}.Routes())
	srv.ListenAndServe()
}
