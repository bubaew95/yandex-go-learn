package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Port     string
	BaseURL  string
	FilePath string
}

func NewConfig() *Config {
	port := flag.String("a", ":8080", "отвечает за адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "", " отвечает за базовый адрес результирующего сокращённого URL")
	filePath := flag.String("f", "data.json", "путь до файла, куда сохраняются данные в формате JSON")

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

	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		*filePath = envFilePath
	}

	return &Config{
		Port:     *port,
		BaseURL:  *baseURL,
		FilePath: *filePath,
	}
}
