package config

import (
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

	if *serverAddress != "" {
		config.ServerAddress = *serverAddress
	}
	if *baseURL != "" {
		config.BaseURL = *baseURL
	}
	if *filePath != "" {
		config.FilePath = *filePath
	}
	if *databaseDSN != "" {
		config.DataBaseDSN = *databaseDSN
	}
	if *enableHTTPS {
		config.EnableHTTPS = *enableHTTPS
	}

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		config.ServerAddress = envServerAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.BaseURL = envBaseURL
	}

	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		config.FilePath = envFilePath
	}

	if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
		config.DataBaseDSN = envDataBaseDSN
	}

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		enbHTTPS, err := strconv.ParseBool(envEnableHTTPS)
		if err == nil {
			config.EnableHTTPS = enbHTTPS
		}
	}

	if config.ServerAddress == "" {
		config.ServerAddress = "8080"
	}

	if !strings.Contains(config.ServerAddress, ":") {
		config.ServerAddress = ":" + config.ServerAddress
	}

	if config.BaseURL == "" {
		config.BaseURL = fmt.Sprintf("http://localhost%s", config.ServerAddress)
	}

	if config.FilePath == "" {
		config.FilePath = "data.json"
	}

	fmt.Println(config)
	return &config
}

func getFileConfigs(filePath string, cfg *Config) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cfg)
}
