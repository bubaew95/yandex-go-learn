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

	if envBaseUrl := os.Getenv("BASE_URL"); envBaseUrl != "" {
		*baseURL = envBaseUrl
	}

	if *baseURL == "" {
		*baseURL = fmt.Sprintf("http://localhost%s", *port)
	}

	return &Config{
		Port:    *port,
		BaseURL: *baseURL,
	}
}

func NewTestConfig(args []string) *Config {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	port := fs.String("a", ":8080", "Адрес запуска HTTP-сервера")
	baseURL := fs.String("b", "http://localhost", "Базовый адрес результирующего сокращённого URL")

	_ = fs.Parse(args)

	return &Config{
		Port:    *port,
		BaseURL: *baseURL,
	}
}
