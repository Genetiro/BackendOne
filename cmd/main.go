package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Genetiro/BackendOne/internal/server"
	"github.com/Genetiro/BackendOne/internal/transport"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Shortner API
// @version 1.0
// @description API server for links

// @host localhost:8080
// @BasePath /

func main() {
	r := chi.NewRouter()
	srv := server.NewSrv(":8080", r)

	srv.Start()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome")
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	r.Mount("/links", transport.LinkResources{}.Routes())
	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-osSigChan
	log.Println("received OS interrupting signal")
	srv.Stop()

}
