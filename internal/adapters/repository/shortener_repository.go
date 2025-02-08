package repository

import (
	"context"
	"fmt"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

type ShortenerRepository struct {
	shortenerDB storage.ShortenerDB
	cache       map[string]string
}

func NewShortenerRepository(s storage.ShortenerDB) *ShortenerRepository {
	data, _ := s.Load()
	return &ShortenerRepository{
		shortenerDB: s,
		cache:       data,
	}
}

func (s ShortenerRepository) Close() error {
	return s.shortenerDB.Close()
}

func (s ShortenerRepository) SetURL(ctx context.Context, id string, url string) {
	s.cache[id] = url

	data := &model.ShortenURL{
		UUID:        len(s.cache),
		ShortURL:    id,
		OriginalURL: url,
	}

	err := s.shortenerDB.Save(data)
	if err != nil {
		logger.Log.Debug(fmt.Sprintf("Не удалось записать данные в файл. Ошибка: %s", err.Error()))
	}
}

func (s ShortenerRepository) GetURLByID(ctx context.Context, id string) (string, bool) {
	url, ok := s.cache[id]

	return url, ok
}

func (s ShortenerRepository) GetAllURL(ctx context.Context) map[string]string {
	return s.cache
}

func (s ShortenerRepository) Ping() error {
	return nil
}

func (s ShortenerRepository) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error) {
	var responseURLs []model.ShortenerURLResponse

	for _, v := range urls {
		if isEmpty(v.CorrelationID) || isEmpty(v.OriginalURL) {
			continue
		}

		_, existsURL := s.GetURLByID(ctx, v.CorrelationID)
		if existsURL {
			continue
		}

		s.SetURL(ctx, v.CorrelationID, v.OriginalURL)

		responseURLs = append(responseURLs, model.ShortenerURLResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      v.OriginalURL,
		})
	}

	return responseURLs, nil
}
