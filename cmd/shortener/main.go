package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/handlers"
	"github.com/bubaew95/yandex-go-learn/internal/interfaces"
	"github.com/bubaew95/yandex-go-learn/internal/logger"
	"github.com/bubaew95/yandex-go-learn/internal/middlewares"
	"github.com/bubaew95/yandex-go-learn/internal/repository"
	"github.com/bubaew95/yandex-go-learn/internal/service"
	"github.com/bubaew95/yandex-go-learn/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type closer interface {
	Close() error
}

func main() {
	if err := runApp(); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Ошибка запуска приложения: %v", err))
		os.Exit(1)
	}
}

func runApp() error {
	if err := logger.Initialize(); err != nil {
		return fmt.Errorf("ошибка инициализации логирования: %w", err)
	}

	cfg := config.NewConfig()
	shortenerRepository, err := initRepository(*cfg)
	if err != nil {
		return err
	}
	defer safeClose(shortenerRepository)

	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)
	shortenerHandler := handlers.NewShortenerHandler(shortenerService)

	route := setupRouter(shortenerHandler)

	logger.Log.Info("Running server", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(cfg.Port, route); err != nil {
		return fmt.Errorf("ошибка запуска сервера: %w", err)
	}

	return nil
}

func initRepository(cfg config.Config) (interfaces.ShortenerRepositoryInterface, error) {
	fmt.Println(cfg.DataBaseDSN)
	if cfg.DataBaseDSN != "" {
		shortenerRepository, err := repository.NewPgRepository(cfg)
		if err != nil {
			logger.Log.Fatal(fmt.Sprintf("Ошибка инициализации базы данных: %v", err))
		}

		return shortenerRepository, nil
	}

	shortenerDB, err := storage.NewShortenerDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации файла базы данных: %w", err)
	}
	return repository.NewShortenerRepository(*shortenerDB), nil
}

func setupRouter(shortenerHandler *handlers.ShortenerHandler) *chi.Mux {
	route := chi.NewRouter()
	route.Use(middlewares.LoggerMiddleware)
	route.Use(middlewares.GZipMiddleware)

	route.Post("/", shortenerHandler.CreateURL)
	route.Get("/{id}", shortenerHandler.GetURL)
	route.Post("/api/shorten", shortenerHandler.AddNewURL)
	route.Get("/ping", shortenerHandler.Ping)

	return route
}

func safeClose(c closer) {
	if err := c.Close(); err != nil {
		logger.Log.Error(fmt.Sprintf("Ошибка при закрытии ресурса: %v", err))
	}
}
