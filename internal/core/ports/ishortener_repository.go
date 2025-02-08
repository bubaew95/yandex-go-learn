package ports

import (
	"context"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

type ShortenerRepositoryInterface interface {
	GetURLByID(ctx context.Context, id string) (string, bool)
	SetURL(ctx context.Context, id string, url string)
	GetAllURL(ctx context.Context) map[string]string
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error)

	Ping() error
	Close() error
}
