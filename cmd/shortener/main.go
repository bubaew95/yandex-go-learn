package main

import (
	"context"
	"fmt"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/app"
	fileStorage "github.com/bubaew95/yandex-go-learn/internal/adapters/repository/filestorage"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/repository/postgres"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"go.uber.org/zap"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/handlers/http"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
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
	shortenerRepository, err := initRepositoryHelper(*cfg)
	if err != nil {
		return err
	}
	defer safeClose(shortenerRepository)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)
	shortenerService.Run(ctx, &wg)

	shortenerHandler := handlers.NewShortenerHandler(shortenerService, *cfg)

	app := app.NewApp(cfg, *shortenerHandler, shortenerRepository)
	app.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-ch

	app.Stop()

	shortenerService.Close()
	wg.Wait()
	return nil
}

func initRepositoryHelper(cfg config.Config) (service.ShortenerRepository, error) {
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

func safeClose(c closer) {
	if err := c.Close(); err != nil {
		logger.Log.Error("Error when closing a resource", zap.Error(err))
	}
}
