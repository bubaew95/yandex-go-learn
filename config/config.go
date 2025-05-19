package config

import (
	"flag"
	"fmt"
	"os"
)

// Config содержит параметры конфигурации приложения.
type Config struct {
	// Port на котором запускается web сервер
	Port string

	//BaseURL базовый адрес
	BaseURL string

	// FilePath путь к файлу
	FilePath string

	// DataBaseDSN строка подключения к базе данных
	DataBaseDSN string

	// EnableHTTPS Включить https протокол
	EnableHTTPS string
}

// NewConfig создает и возвращает структуру конфигурации Config,
// комбинируя значения из флагов командной строки и переменных окружения.
func NewConfig() *Config {
	port := flag.String("a", ":8080", "отвечает за адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "", " отвечает за базовый адрес результирующего сокращённого URL")
	filePath := flag.String("f", "data.json", "путь до файла, куда сохраняются данные в формате JSON")
	databaseDSN := flag.String("d", "", "Строка подключения к БД")
	enableHTTPS := flag.String("s", "", "Включить https протокол")

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

	if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
		*databaseDSN = envDataBaseDSN
	}

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		*enableHTTPS = envEnableHTTPS
	}

	return &Config{
		Port:        *port,
		BaseURL:     *baseURL,
		FilePath:    *filePath,
		DataBaseDSN: *databaseDSN,
		EnableHTTPS: *enableHTTPS,
	}
}
