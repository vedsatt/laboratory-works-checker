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
	port string
}

// Создает новый объект сервера
func New(port string) *Server {
	s := &Server{
		port: fmt.Sprintf(":%s", port),
	}

	return s
}

// Запускает сервер без поддержки SSL/TLS
// Функция используется для тестирования при разработке
// Перед выводом сервера в продакшн следует получить SSL/TLS сертификат и
// использовать надежное подключение
func (s *Server) Start() {
	submit := http.HandlerFunc(SubmitLabHandler)

	// Мьюх нужен, чтобы обслуживать мидлвееры
	mux := http.NewServeMux()

	// Хендлер для отправки работ на проверку
	mux.Handle("/api/v1/submit", logsMiddleware(enableCORS(submit)))

	log.Printf("Starting HTTP server on port: %v", s.port)
	log.Fatal(http.ListenAndServe(s.port, mux))
}
