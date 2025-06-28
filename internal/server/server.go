package server

import (
	"context"
	"log"
	"net/http"
)

type HttpServer struct {
	srv *http.Server
}

func NewServer(addr string, handler http.Handler) *HttpServer {
	return &HttpServer{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *HttpServer) Start() error {
	log.Printf("Starting server on %s", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	log.Println("Stopping server...")
	return s.srv.Shutdown(ctx)
}
