package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	srv http.Server
}

func NewSrv(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return s
}

func (s *Server) Start() {
	go s.srv.ListenAndServe()
	log.Println("server started on", s.srv.Addr)
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	s.srv.Shutdown(ctx)
	log.Println("server turned off")
}
