// Package server запускает http/https сервер, который принимает лаботараторные работы на проверку
package server

import (
	"fmt"
	"log"
	"net/http"
)

// Структура сервера. С помощью нее создается объект сервера с методом Start()
// Порт указывается в переменной среде, в случае отсутствия - берется значенеи по умолчанию (80)
type Server struct {
	port      string
	httpsPort string
	certFile  string
	keyFile   string
}

// Создает новый объект сервера
func New(port, httpsPort, certFile, keyFile string) *Server {
	s := &Server{
		port:      fmt.Sprintf(":%s", port),
		httpsPort: fmt.Sprintf(":%s", httpsPort),
		certFile:  certFile,
		keyFile:   keyFile,
	}

	return s
}

// Запускает сервер с поддержкой SSL протокола
// certFile - путь к SSL-сертификату (например, "cert.pem")
// keyFile - путь к приватному ключу (например, "key.pem")
func (s *Server) Start() {
	submit := http.HandlerFunc(SubmitLabHandler)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/submit", logsMiddleware(enableCORS(submit)))

	log.Printf("Starting SSL (HTTPS) server on port: %v", s.port)
	err := http.ListenAndServeTLS(s.port, s.certFile, s.keyFile, mux)
	if err != nil {
		log.Fatalf("Failed to start SSL server: %v", err)
	}
}
