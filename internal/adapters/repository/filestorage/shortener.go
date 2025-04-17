package filestorage

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

type ShortenerRepository struct {
	shortenerDB storage.ShortenerDB
	mx          *sync.RWMutex
	cache       map[string]string
}

func NewShortenerRepository(s storage.ShortenerDB) (*ShortenerRepository, error) {
	data, err := s.Load()
	if err != nil {
		return nil, err
	}

	return &ShortenerRepository{
		shortenerDB: s,
		mx:          &sync.RWMutex{},
		cache:       data,
	}, nil
}

func (s ShortenerRepository) Close() error {
	return s.shortenerDB.Close()
}

func (s ShortenerRepository) SetURL(ctx context.Context, id string, url string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.cache[id] = url

	data := &model.ShortenURL{
		UUID:        len(s.cache),
		ShortURL:    id,
		OriginalURL: url,
	}

	return s.shortenerDB.Save(data)
}

func (s ShortenerRepository) GetURLByID(ctx context.Context, id string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	url, ok := s.cache[id]
	if !ok {
		return "", errors.New("not found")
	}

	return url, nil
}

func (s ShortenerRepository) GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for id, v := range s.cache {
		if strings.Contains(v, originalURL) {
			return id, true
		}
	}

	return "", false
}

func (s ShortenerRepository) Ping(ctx context.Context) error {
	return nil
}

func (s ShortenerRepository) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error {
	for _, v := range urls {
		_, ok := s.cache[v.CorrelationID]
		if !ok {
			continue
		}

		s.SetURL(ctx, v.CorrelationID, v.OriginalURL)
	}

	return nil
}

func (s ShortenerRepository) GetURLSByUserID(ctx context.Context, userID string) (map[string]string, error) {
	return nil, nil
}

func (s ShortenerRepository) DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error {
	return nil
}
