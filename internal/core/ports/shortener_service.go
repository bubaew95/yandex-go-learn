package ports

import (
	"context"
	"sync"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

type ShortenerService interface {
	GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error)
	GetURLByID(ctx context.Context, id string) (string, bool)
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)
	GetAllURL(ctx context.Context) map[string]string
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error)
	GetURLSByUserID(ctx context.Context, userID string) ([]model.ShortenerURLSForUserResponse, error)
	DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error

	Worker(ctx context.Context, wg *sync.WaitGroup)
	ScheduleURLDeletion(ctx context.Context, items []model.URLToDelete)

	RandStringBytes(n int) string
	Ping() error
}
