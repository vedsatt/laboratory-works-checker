// Package config получает данные из переменных сред.
package config

import (
	"log"
	"os"
)

var (
	defPort string = "80"
)

// Конфиг сервера
type Config struct {
	Port string
}

// Возвращает порт сервера из переменных сред
func GetConfig() *Config {
	// Устанавливаем значения по умолчанию сразу
	c := &Config{
		Port: defPort,
	}

	if port, ok := os.LookupEnv("PORT"); !ok {
		log.Printf("Port is unset. Using default value: %v", defPort)
	} else {
		c.Port = port
	}

	return c
}
