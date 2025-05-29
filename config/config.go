package config

import (
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

// Config содержит параметры конфигурации приложения.
type Config struct {
	// ServerAddress на котором запускается web сервер
	ServerAddress string `json:"server_address"`

	//BaseURL базовый адрес
	BaseURL string `json:"base_path"`

	// FilePath путь к файлу
	FilePath string `json:"file_storage_path"`

	// DataBaseDSN строка подключения к базе данных
	DataBaseDSN string `json:"database_dsn"`

	// EnableHTTPS Включить https протокол
	EnableHTTPS bool `json:"enable_https"`
}

// NewConfig создает и возвращает структуру конфигурации Config,
// комбинируя значения из флагов командной строки и переменных окружения.
func NewConfig() *Config {

	var config Config
	var fileConfigPath string

	serverAddress := flag.String("a", "", "Адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "", "Базовый адрес сокращённого URL")
	filePath := flag.String("f", "", "Путь до JSON-файла")
	databaseDSN := flag.String("d", "", "Строка подключения к базе данных")
	enableHTTPS := flag.Bool("s", false, "Включить HTTPS")

	flag.StringVar(&fileConfigPath, "c", "", "Путь к JSON файлу конфигурации")
	flag.StringVar(&fileConfigPath, "config", "", "Путь к JSON файлу конфигурации")

	flag.Parse()

	if envFileConfig := os.Getenv("CONFIG"); envFileConfig != "" {
		fileConfigPath = envFileConfig
	}

	if fileConfigPath != "" {
		err := getFileConfigs(fileConfigPath, &config)
		if err != nil {
			logger.Log.Debug("Config file error", zap.Error(err))
		}
	}

	config.ServerAddress = cmp.Or(os.Getenv("SERVER_ADDRESS"), *serverAddress, config.ServerAddress, "8080")

	if !strings.Contains(config.ServerAddress, ":") {
		config.ServerAddress = ":" + config.ServerAddress
	}

	config.BaseURL = cmp.Or(os.Getenv("BASE_URL"), *baseURL, config.BaseURL, fmt.Sprintf("http://localhost%s", config.ServerAddress))
	config.FilePath = cmp.Or(os.Getenv("FILE_STORAGE_PATH"), *filePath, config.FilePath, "data.json")
	config.DataBaseDSN = cmp.Or(os.Getenv("DATABASE_DSN"), *databaseDSN, config.DataBaseDSN)

	if *enableHTTPS {
		config.EnableHTTPS = *enableHTTPS
	}

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		enbHTTPS, err := strconv.ParseBool(envEnableHTTPS)
		if err == nil {
			config.EnableHTTPS = enbHTTPS
		}
	}

	return &config
}

func getFileConfigs(filePath string, cfg *Config) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cfg)
}
