package server

import "net/http"

func logsMiddleware(next http.HandleFunc())