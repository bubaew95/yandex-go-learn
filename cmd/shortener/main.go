package main

import (
	"net/http"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/handlers"
	"github.com/bubaew95/yandex-go-learn/internal/logger"
	"github.com/bubaew95/yandex-go-learn/internal/middlewares"
	"github.com/bubaew95/yandex-go-learn/internal/repository"
	"github.com/bubaew95/yandex-go-learn/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}

	data := make(map[string]string)

	cfg := config.NewConfig()
	shortenerRepository := repository.NewShortenerRepository(data, cfg.BaseURL)
	shortenerService := service.NewShortenerService(shortenerRepository)
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
