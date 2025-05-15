package main

import (
	"github.com/vedsatt/laboratory-works-checker/internal/config"
	server "github.com/vedsatt/laboratory-works-checker/internal/http-server"
	tcpclient "github.com/vedsatt/laboratory-works-checker/internal/tcp-client"
)

func main() {
	cfg := config.GetConfig()

	server := server.New(cfg.Port)

	go server.Start()

	tcpclient.Connect()
}
