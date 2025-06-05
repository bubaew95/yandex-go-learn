package helper

import (
	"fmt"
	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	fileStorage "github.com/bubaew95/yandex-go-learn/internal/adapters/repository/filestorage"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/repository/postgres"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/service"
	"go.uber.org/zap"
)

func InitRepositoryHelper(cfg config.Config) (service.ShortenerRepository, error) {
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

func SafeClose(c closer) {
	if err := c.Close(); err != nil {
		logger.Log.Error("Error when closing a resource", zap.Error(err))
	}
}
