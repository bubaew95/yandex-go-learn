package main

import (
	"context"
	"errors"
	"fmt"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/handlers"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/handlers/middleware"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	fileStorage "github.com/bubaew95/yandex-go-learn/internal/adapters/repository/filestorage"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/repository/postgres"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/service"
)

type closer interface {
	Close() error
}

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

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
	shortenerService.Run(ctx, &wg)

	shortenerHandler := handlers.NewShortenerHandler(shortenerService)
	route := setupRouter(shortenerHandler)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: route,
	}

	if cfg.EnableHTTPS {
		logger.Log.Info("Running https server", zap.String("port", cfg.ServerAddress))
		go func() {
			if err := server.Serve(autocert.NewListener("example.com")); err != nil {
				logger.Log.Error("Failed to start https(tsl) server", zap.Error(err))
			}
		}()
	} else {
		logger.Log.Info("Running server", zap.String("ports", cfg.ServerAddress))
		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Log.Error("Failed to start http server", zap.Error(err))
			}
		}()
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-ch

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	logger.Log.Info("Shutting down...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Info("Http server shutdown error", zap.Error(err))
	}

	shortenerService.Close()

	wg.Wait()
	return nil
}

func initRepository(cfg config.Config) (service.ShortenerRepository, error) {
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

	route.Mount("/debug", chi_middleware.Profiler())

	return route
}

func safeClose(c closer) {
	if err := c.Close(); err != nil {
		logger.Log.Error("Error when closing a resource", zap.Error(err))
	}
}
