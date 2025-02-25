package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/handlers"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/handlers/middleware"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	fileStorage "github.com/bubaew95/yandex-go-learn/internal/adapters/repository/filestorage"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/repository/postgres"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
	"github.com/bubaew95/yandex-go-learn/internal/core/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type closer interface {
	Close() error
}

func main() {
	if err := runApp(); err != nil {
		logger.Log.Fatal("Application startup error", zap.Error(err))
	}
}

func runApp() error {
	if err := logger.Initialize(); err != nil {
		return fmt.Errorf("logging initialization error: %w", err)
	}

	cfg := config.NewConfig()
	shortenerRepository, err := initRepository(*cfg)
	if err != nil {
		return err
	}
	defer safeClose(shortenerRepository)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)
	shortenerService.Worker(ctx, &wg)

	shortenerHandler := handlers.NewShortenerHandler(shortenerService)
	route := setupRouter(shortenerHandler)

	logger.Log.Info("Running server", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(cfg.Port, route); err != nil {
		return fmt.Errorf("server startup error: %w", err)
	}

	return nil
}

func initRepository(cfg config.Config) (ports.ShortenerRepository, error) {
	if cfg.DataBaseDSN != "" {
		shortenerRepository, err := postgres.NewShortenerRepository(cfg)
		if err != nil {
			logger.Log.Fatal("Database initialization error", zap.Error(err))
		}

		return shortenerRepository, nil
	}

	shortenerDB, err := storage.NewShortenerDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("database file initialization error: %w", err)
	}

	shortener, err := fileStorage.NewShortenerRepository(*shortenerDB)

	return shortener, err
}

func setupRouter(shortenerHandler *handlers.ShortenerHandler) *chi.Mux {
	route := chi.NewRouter()
	route.Use(middleware.LoggerMiddleware)
	route.Use(middleware.GZipMiddleware)
	route.Use(middleware.CookieMiddleware)

	route.Post("/", shortenerHandler.CreateURL)
	route.Get("/{id}", shortenerHandler.GetURL)
	route.Get("/ping", shortenerHandler.Ping)

	route.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", shortenerHandler.AddNewURL)
		r.Post("/batch", shortenerHandler.Batch)
	})

	route.Route("/api/user", func(r chi.Router) {
		r.Get("/urls", shortenerHandler.GetUserURLS)
		r.Delete("/urls", shortenerHandler.DeleteUserURLS)
	})

	return route
}

func safeClose(c closer) {
	if err := c.Close(); err != nil {
		logger.Log.Error("Error when closing a resource", zap.Error(err))
	}
}
