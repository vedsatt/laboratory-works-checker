package server

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port string
}

// Создает новый объект сервера
func New(port string) *Server {
	s := &Server{
		port: fmt.Sprintf(":%s", port),
	}

	return s
}

// Запуск сервера
func (s *Server) Start() {
	submit := http.HandlerFunc(SubmitLabHandler)

	// Мьюх нужен, чтобы обслуживать мидлвееры
	mux := http.NewServeMux()

	// Хендлер для отправки работ на проверку
	mux.Handle("/api/v1/submit", logsMiddleware(enableCORS(submit)))

	log.Printf("Starting server on port: %v", s.port)
	http.ListenAndServe(s.port, mux)
}
