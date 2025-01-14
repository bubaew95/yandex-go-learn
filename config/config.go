package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Port    string
	BaseURL string
}

func NewConfig() *Config {
	port := flag.String("a", ":8080", "отвечает за адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "", " отвечает за базовый адрес результирующего сокращённого URL")

	flag.Parse()

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		*port = envServerAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		*baseURL = envBaseURL
	}

	if *baseURL == "" {
		*baseURL = fmt.Sprintf("http://localhost%s", *port)
	}

	return &Config{
		Port:    *port,
		BaseURL: *baseURL,
	}
}
