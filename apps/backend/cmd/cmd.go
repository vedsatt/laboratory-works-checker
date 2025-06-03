package main

import (
	"github.com/vedsatt/laboratory-works-checker/apps/backend/internal/cleaner"
	"github.com/vedsatt/laboratory-works-checker/apps/backend/internal/config"
	server "github.com/vedsatt/laboratory-works-checker/apps/backend/internal/http-server"
)

func main() {
	// Очищаем директорию от мусорных папок
	cleaner.Clean()

	// Получаем данные из переменных сред (еще одно удобство докера)
	cfg := config.GetConfig()

	// создаем объект сервера, настраиваем и запускаем его
	server := server.New(cfg.Port, cfg.HttpsPort, "server.crt", "server.key")
	server.Start()
}
