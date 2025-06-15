// Package config получает данные из переменных сред.
package config

import (
	"log"
	"os"
)

var (
	defPort      string = "80"
	defHttpsPort string = "8443"
)

// Конфиг сервера
type Config struct {
	Port      string
	HttpsPort string
}

// Возвращает порт сервера из переменных сред
func GetConfig() *Config {
	// Устанавливаем значения по умолчанию сразу
	c := &Config{
		Port:      defPort,
		HttpsPort: defHttpsPort,
	}

	if port, ok := os.LookupEnv("PORT"); !ok {
		log.Printf("Port is unset. Using default value: %v", defPort)
	} else {
		c.Port = port
	}

	if httpsPort, ok := os.LookupEnv("HTTPS_PORT"); !ok {
		log.Printf("HTTPS port is unset. Using default value: %v", defHttpsPort)
	} else {
		c.HttpsPort = httpsPort
	}

	return c
}
