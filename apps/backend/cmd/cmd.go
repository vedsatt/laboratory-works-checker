// Package main звпускает cleaner, получает конфиг, настраивает и запускает сервер
package main

import (
	"github.com/vedsatt/laboratory-works-checker/apps/backend/application/cleaner"
	"github.com/vedsatt/laboratory-works-checker/apps/backend/application/config"
	server "github.com/vedsatt/laboratory-works-checker/apps/backend/application/http-server"
)

func main() {
	// Очищаем директорию от мусорных папок
	cleaner.Clean()

	// Получаем данные из переменных сред (еще одно удобство докера)
	cfg := config.GetConfig()

	// создает объект сервера и настраивает его
	// cert.pem - сертификат TLS
	// key.pem - ключ
	// этих файлов нет и они служат заглушкой при тестировании
	server := server.New(cfg.Port, cfg.HttpsPort, "cert.pem", "key.pem")
	server.Start()
}
