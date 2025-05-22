package main

import (
	"github.com/vedsatt/laboratory-works-checker/apps/backend/internal/config"
	server "github.com/vedsatt/laboratory-works-checker/apps/backend/internal/http-server"
)

func main() {
	cfg := config.GetConfig()

	server := server.New(cfg.Port)

	server.Start()
}
