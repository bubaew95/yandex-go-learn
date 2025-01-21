package repository

import (
	"fmt"

	"github.com/bubaew95/yandex-go-learn/internal/logger"
	"github.com/bubaew95/yandex-go-learn/internal/models"
	"github.com/bubaew95/yandex-go-learn/internal/tools"
)

type ShortenerRepository struct {
	shortenerDB tools.ShortenerDB
}

func NewShortenerRepository(s tools.ShortenerDB) *ShortenerRepository {
	return &ShortenerRepository{
		shortenerDB: s,
	}
}

func (s ShortenerRepository) SetURL(id string, url string) {
	data := &models.ShortenURL{
		UUID:        s.shortenerDB.Count() + 1,
		ShortURL:    id,
		OriginalURL: url,
	}

	err := s.shortenerDB.Save(data)
	if err != nil {
		logger.Log.Debug(fmt.Sprintf("Не удалось записать данные в файл. Ошибка: %s", err.Error()))
	}
}

func (s ShortenerRepository) GetURLByID(id string) (string, bool) {
	data, err := s.shortenerDB.Load()
	if err != nil {
		return "", false
	}

	url, ok := data[id]
	return url, ok
}

func (s ShortenerRepository) GetAllURL() map[string]string {
	data, _ := s.shortenerDB.Load()

	return data
}
