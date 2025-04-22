// Package filestorage предоставляет реализацию репозитория сокращённых URL,
// основанную на файловом хранилище с поддержкой кэширования в памяти.
package filestorage

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

// ShortenerRepository реализует интерфейс репозитория для работы с сокращёнными URL.
// Использует in-memory кэш с синхронизацией и файловое хранилище.
type ShortenerRepository struct {
	shortenerDB storage.ShortenerDB
	mx          *sync.RWMutex
	cache       map[string]string
}

// NewShortenerRepository инициализирует новый экземпляр ShortenerRepository.
// Загружает данные из хранилища в кэш.
//
// Возвращает ошибку, если загрузка данных не удалась.
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

// Close закрывает соединение с хранилищем.
func (s ShortenerRepository) Close() error {
	return s.shortenerDB.Close()
}

// SetURL сохраняет соответствие между коротким ID и оригинальным URL.
// Добавляет запись в кэш и в файловое хранилище.
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

// GetURLByID возвращает оригинальный URL по его короткому идентификатору.
// Возвращает ошибку, если соответствие не найдено.
func (s ShortenerRepository) GetURLByID(ctx context.Context, id string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	url, ok := s.cache[id]
	if !ok {
		return "", errors.New("not found")
	}

	return url, nil
}

// GetURLByOriginalURL возвращает короткий ID по оригинальному URL.
// Используется сравнение по вхождению (contains).
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

// Ping реализует метод "пинга" для проверки доступности хранилища.
// В текущей реализации всегда возвращает nil.
func (s ShortenerRepository) Ping(ctx context.Context) error {
	return nil
}

// InsertURLs добавляет список URL в хранилище, если они ещё не существуют.
// Пропускает уже существующие записи.
func (s ShortenerRepository) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error {
	for _, v := range urls {
		_, err := s.GetURLByID(ctx, v.CorrelationID)
		if err == nil {
			continue
		}

		err = s.SetURL(ctx, v.CorrelationID, v.OriginalURL)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertURLTwo обновляет только те записи, которые уже есть в кэше.
// Полезно для обновления оригинальных URL.
func (s ShortenerRepository) InsertURLTwo(ctx context.Context, urls []model.ShortenerURLMapping) error {
	for _, v := range urls {
		_, ok := s.cache[v.CorrelationID]
		if !ok {
			continue
		}

		err := s.SetURL(ctx, v.CorrelationID, v.OriginalURL)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetURLSByUserID возвращает список URL, привязанных к конкретному пользователю.
// В текущей реализации не реализован и всегда возвращает nil.
func (s ShortenerRepository) GetURLSByUserID(ctx context.Context, userID string) (map[string]string, error) {
	return nil, nil
}

// DeleteUserURLS удаляет список URL, привязанных к пользователю.
// В текущей реализации не реализован и всегда возвращает nil.
func (s ShortenerRepository) DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error {
	return nil
}
