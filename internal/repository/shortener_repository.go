package repository

import (
	"fmt"

	"github.com/bubaew95/yandex-go-learn/internal/logger"
	"github.com/bubaew95/yandex-go-learn/internal/models"
	"github.com/bubaew95/yandex-go-learn/internal/storage"
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

func (s ShortenerRepository) SetURL(id string, url string) {
	s.cache[id] = url

	data := &models.ShortenURL{
		UUID:        len(s.cache),
		ShortURL:    id,
		OriginalURL: url,
	}

	err := s.shortenerDB.Save(data)
	if err != nil {
		logger.Log.Debug(fmt.Sprintf("Не удалось записать данные в файл. Ошибка: %s", err.Error()))
	}
}

func (s ShortenerRepository) GetURLByID(id string) (string, bool) {
	url, ok := s.cache[id]

	return url, ok
}

func (s ShortenerRepository) GetAllURL() map[string]string {
	return s.cache
}

func (s ShortenerRepository) Ping() error {
	return nil
}
