package ports

import (
	"context"
	"errors"
	"sync"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

var (
	ErrUniqueIndex = errors.New("url already exists")
)

type ShortenerRepository interface {
	GetURLByID(ctx context.Context, id string) (string, bool)
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)
	SetURL(ctx context.Context, id string, url string) error
	GetAllURL(ctx context.Context) map[string]string
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error
	GetURLSByUserID(ctx context.Context, userID string) (map[string]string, error)
	DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error

	Ping() error
	Close() error
}

type ShortenerService interface {
	GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error)
	GetURLByID(ctx context.Context, id string) (string, bool)
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)
	GetAllURL(ctx context.Context) map[string]string
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error)
	GetURLSByUserID(ctx context.Context, userID string) ([]model.ShortenerURLSForUserResponse, error)
	DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error

	Run(ctx context.Context, wg *sync.WaitGroup)
	ScheduleURLDeletion(ctx context.Context, items []model.URLToDelete)

	RandStringBytes(n int) string
	Ping() error
}
