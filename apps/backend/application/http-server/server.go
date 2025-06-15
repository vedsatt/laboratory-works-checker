// Package server запускает http/https сервер, который принимает лаботараторные работы на проверку
package server

import (
	"crypto/tls"
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

// Запускает сервер с поддержкой TLS протокола
func (s *Server) Start() {
	submit := http.HandlerFunc(SubmitLabHandler)

	// Мьюх нужен, чтобы обслуживать мидлвееры
	mux := http.NewServeMux()

	// Хендлер для отправки работ на проверку
	mux.Handle("/api/v1/submit", logsMiddleware(enableCORS(submit)))

	go func() {
		log.Printf("Starting HTTP server on port: %v (redirect to HTTPS)", s.port)
		err := http.ListenAndServe(s.port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+s.httpsPort+r.RequestURI, http.StatusMovedPermanently)
		}))
		if err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Настройка TLS конфигурации
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// Настройка HTTPS сервера
	server := &http.Server{
		Addr:         s.httpsPort,
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.Printf("Starting HTTPS server on port: %v", s.httpsPort)
	err := server.ListenAndServeTLS(s.certFile, s.keyFile)
	if err != nil {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}
