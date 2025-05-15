package config

import (
	"log"
	"os"
)

var (
	defaultPort string = "80"
	defaultURL  string = "backend.labchecker"
)

type Config struct {
	Port string
	URL  string
}

func GetConfig() *Config {
	c := &Config{
		Port: defaultPort,
		URL:  defaultURL,
	}

	if port, ok := os.LookupEnv("PORT"); !ok && port != "" {
		log.Printf("Port is unset. Using default value: %v", defaultPort)
	} else {
		c.Port = port
	}

	if url, ok := os.LookupEnv("SERVER_URL"); !ok && url != "" {
		log.Printf("Server url is unset. Using default value: %v", defaultURL)
	} else {
		c.URL = url
	}

	return c
}
