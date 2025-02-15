package ports

import (
	"context"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

type ShortenerService interface {
	GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error)
	GetURLByID(ctx context.Context, id string) (string, bool)
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)
	GetAllURL(ctx context.Context) map[string]string
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error)

	RandStringBytes(n int) string
	Ping() error
}
