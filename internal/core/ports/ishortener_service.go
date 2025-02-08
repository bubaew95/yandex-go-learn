package ports

import (
	"context"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

type ShortenerServiceInterface interface {
	GenerateURL(ctx context.Context, url string, randomStringLength int) string
	GetURLByID(ctx context.Context, id string) (string, bool)
	GetAllURL(ctx context.Context) map[string]string
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error)

	RandStringBytes(n int) string
	Ping() error
}
