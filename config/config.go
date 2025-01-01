package config

import (
	"flag"
)

type Config struct {
	Port    string
	BaseURL string
}

func NewConfig() *Config {
	port := flag.String("a", "8080", "отвечает за адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "http://localhost", " отвечает за базовый адрес результирующего сокращённого URL")

	flag.Parse()

	return &Config{
		Port:    *port,
		BaseURL: *baseURL,
	}
}

func NewTestConfig(args []string) *Config {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	port := fs.String("a", "8080", "Адрес запуска HTTP-сервера")
	baseURL := fs.String("b", "http://localhost", "Базовый адрес результирующего сокращённого URL")

	_ = fs.Parse(args)

	return &Config{
		Port:    *port,
		BaseURL: *baseURL,
	}
}
