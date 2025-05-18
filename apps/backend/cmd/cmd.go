package main

import (
	"github.com/vedsatt/laboratory-works-checker/internal/config"
	server "github.com/vedsatt/laboratory-works-checker/internal/http-server"
)

func main() {
	cfg := config.GetConfig()

	server := server.New(cfg.Port)

	server.Start()
}
