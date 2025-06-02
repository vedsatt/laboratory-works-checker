// Package server является http-сервером, через который ппроисходит взаиможействие с checker
package server

import (
	"fmt"
	"log"
	"net/http"
)

// Структура сервера. С помощью нее создается объект сервера с методом Start()
// Порт указывается в переменной среде, в случае отсутствия - берется значенеи по умолчанию (80)
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

// Запускает сервер
func (s *Server) Start() {
	submit := http.HandlerFunc(SubmitLabHandler)

	// Мьюх нужен, чтобы обслуживать мидлвееры
	mux := http.NewServeMux()

	// Хендлер для отправки работ на проверку
	mux.Handle("/api/v1/submit", logsMiddleware(enableCORS(submit)))

	log.Printf("Starting server on port: %v", s.port)
	// http.ListenAndServeTLS(s.port, "server.crt", "server.key", mux)
	http.ListenAndServe(s.port, mux)
}
