package ports

import (
	"context"
	"errors"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

type ShortenerRepositoryInterface interface {
	GetURLByID(ctx context.Context, id string) (string, bool)
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)
	SetURL(ctx context.Context, id string, url string) error
	GetAllURL(ctx context.Context) map[string]string
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error

	Ping() error
	Close() error
}

var (
	ErrUniqueIndex = errors.New("url already exists")
)
