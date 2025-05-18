package server

import (
	"log"
	"net/http"
	"time"
)

func logsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %v, URL: %v", r.Method, r.URL)
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("URL: %s, completion time: %v", r.URL, duration)
	})
}
