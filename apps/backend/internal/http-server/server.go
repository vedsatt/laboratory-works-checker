package server

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port string
}

func New(port string) *Server {
	s := &Server{
		port: fmt.Sprintf(":%s", port),
	}

	return s
}

func (s *Server) Start() {
	submit := http.HandlerFunc(SubmitLabHandler)

	mux := http.NewServeMux()

	mux.Handle("/api/v1/submit", logsMiddleware(enableCORS(submit)))

	log.Printf("Starting server on port: %v", s.port)
	http.ListenAndServe(s.port, mux)
}
