package main

import (
	"fmt"
	"net/http"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/handlers"
	"github.com/bubaew95/yandex-go-learn/internal/logger"
	"github.com/bubaew95/yandex-go-learn/internal/middlewares"
	"github.com/bubaew95/yandex-go-learn/internal/repository"
	"github.com/bubaew95/yandex-go-learn/internal/service"
	"github.com/bubaew95/yandex-go-learn/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}

	cfg := config.NewConfig()
	shortenerDB, err := storage.NewShortenerDB(*cfg)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("Ошибка инициализации базы данных: %v", err))
	}

	defer func() {
		if err := shortenerDB.Close(); err != nil {
			logger.Log.Error(fmt.Sprintf("Ошибка при закрытии базы данных: %v", err))
		}
	}()

	shortenerRepository := repository.NewShortenerRepository(*shortenerDB)
	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)
	shortenerHandler := handlers.NewShortenerHandler(shortenerService)

	route := chi.NewRouter()
	route.Use(middlewares.LoggerMiddleware)
	route.Use(middlewares.GZipMiddleware)

	route.Post("/", shortenerHandler.CreateURL)
	route.Get("/{id}", shortenerHandler.GetURL)
	route.Post("/api/shorten", shortenerHandler.AddNewURL)

	if err := run(cfg, route); err != nil {
		panic(err)
	}
}

func run(cfg *config.Config, route *chi.Mux) error {
	logger.Log.Info("Running server", zap.String("port", cfg.Port))

	return http.ListenAndServe(cfg.Port, route)
}
