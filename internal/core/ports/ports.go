package ports

import (
	"context"
	"sync"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

type ShortenerRepository interface {
	GetURLByID(ctx context.Context, id string) (string, error)
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)
	SetURL(ctx context.Context, id string, url string) error
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error
	GetURLSByUserID(ctx context.Context, userID string) (map[string]string, error)
	DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error

	Ping() error
	Close() error
}

type ShortenerService interface {
	GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error)
	GetURLByID(ctx context.Context, id string) (string, error)
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error)
	GetURLSByUserID(ctx context.Context, userID string) ([]model.ShortenerURLSForUserResponse, error)
	DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error

	Run(ctx context.Context, wg *sync.WaitGroup)
	ScheduleURLDeletion(ctx context.Context, items []model.URLToDelete)

	RandStringBytes(n int) string
	Ping() error
}
