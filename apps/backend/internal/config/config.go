package config

import (
	"log"
	"os"
)

var (
	defPort string = "80"
	defURL  string = "localhost"
)

type Config struct {
	Port string
	URL  string
}

// Получаем порт и юрл сервера из переменных сред
func GetConfig() *Config {
	c := &Config{
		Port: defPort,
		URL:  defURL,
	}

	if port, ok := os.LookupEnv("PORT"); !ok {
		log.Printf("Port is unset. Using default value: %v", defPort)
	} else {
		c.Port = port
	}

	if url, ok := os.LookupEnv("SERVER_URL"); !ok {
		log.Printf("Server url is unset. Using default value: %v", defURL)
	} else {
		c.URL = url
	}

	return c
}
