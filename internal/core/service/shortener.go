// Package service содержит реализацию бизнес-логики сервиса сокращения ссылок,
// включая генерацию уникальных ID, хранение, получение и удаление URL.
package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

// ShortenerRepository определяет контракт для репозитория сокращённых URL.
// Этот интерфейс реализуется различными адаптерами хранилищ (например, файловая система, PostgreSQL).
//
//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=ShortenerRepository --filename=repositoryemock_test.go --inpackage
type ShortenerRepository interface {
	// GetURLByID возвращает оригинальный URL по его сокращённому идентификатору.
	GetURLByID(ctx context.Context, id string) (string, error)

	// GetURLByOriginalURL ищет короткий ID по оригинальному URL.
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)

	// SetURL сохраняет соответствие между коротким ID и оригинальным URL.
	SetURL(ctx context.Context, id string, url string) error

	// InsertURLs добавляет список сокращённых URL (например, при массовом импорте).
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) error

	// GetURLSByUserID возвращает карту всех сокращённых ссылок, привязанных к пользователю.
	GetURLSByUserID(ctx context.Context, userID string) (map[string]string, error)

	// DeleteUserURLS помечает ссылки как удалённые по запросу пользователя.
	DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error

	// Ping проверяет доступность репозитория.
	Ping(ctx context.Context) error

	// Close освобождает ресурсы (например, соединения с БД).
	Close() error

	// Stats возвращает статистику
	Stats(ctx context.Context) (model.StatsRespose, error)
}

// ShortenerService реализует бизнес-логику для сокращения URL.
// Поддерживает генерацию уникальных ссылок, сохранение, извлечение и отложенное удаление.
type ShortenerService struct {
	repository ShortenerRepository
	config     config.Config
	mx         *sync.Mutex
	deleteChan chan model.URLToDelete
}

// NewShortenerService создаёт и инициализирует новый экземпляр ShortenerService.
// Принимает хранилище и конфигурацию приложения.
func NewShortenerService(r ShortenerRepository, cfg config.Config) *ShortenerService {
	return &ShortenerService{
		repository: r,
		config:     cfg,
		mx:         &sync.Mutex{},
		deleteChan: make(chan model.URLToDelete),
	}
}

// GenerateURL генерирует уникальный идентификатор для заданного URL и сохраняет его.
// Повторяет генерацию, пока не будет найден уникальный ID.
func (s ShortenerService) GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	var genID string
	for {
		genID = s.RandStringBytes(randomStringLength)
		_, err := s.repository.GetURLByID(ctx, genID)

		if err != nil {
			err := s.repository.SetURL(ctx, genID, url)
			if err != nil {
				return "", err
			}

			break
		}
	}

	return s.generateResponseURL(genID), nil
}

// RandStringBytes генерирует случайную строку заданной длины из латинских букв.
func (s ShortenerService) RandStringBytes(n int) string {

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

// GetURLByID возвращает оригинальный URL по короткому ID.
func (s ShortenerService) GetURLByID(ctx context.Context, id string) (string, error) {
	return s.repository.GetURLByID(ctx, id)
}

// GetURLByOriginalURL возвращает короткий URL по оригинальному, если он уже существует.
func (s ShortenerService) GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool) {
	id, ok := s.repository.GetURLByOriginalURL(ctx, originalURL)

	if ok {
		return s.generateResponseURL(id), ok
	}

	return id, ok
}

// Ping - проверка подключения в БД
func (s ShortenerService) Ping(ctx context.Context) error {
	return s.repository.Ping(ctx)
}

func (s ShortenerService) generateResponseURL(id string) string {
	return fmt.Sprintf("%s/%s", s.config.BaseURL, id)
}

// InsertURLs сохраняет пакет сокращённых ссылок и возвращает сгенерированные короткие ссылки.
func (s ShortenerService) InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error) {
	var items []model.ShortenerURLMapping
	for _, v := range urls {
		if isEmpty(v.CorrelationID) || isEmpty(v.OriginalURL) {
			continue
		}

		items = append(items, v)
	}

	err := s.repository.InsertURLs(ctx, items)
	if err != nil {
		return nil, err
	}

	var responseURLs []model.ShortenerURLResponse
	for _, v := range items {
		responseURLs = append(responseURLs, model.ShortenerURLResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      s.generateResponseURL(v.CorrelationID),
		})
	}

	return responseURLs, nil
}

func isEmpty(t string) bool {
	return strings.TrimSpace(t) == ""
}

// GetURLSByUserID возвращает список сокращённых ссылок, созданных конкретным пользователем.
func (s ShortenerService) GetURLSByUserID(ctx context.Context, userID string) ([]model.ShortenerURLSForUserResponse, error) {
	items, err := s.repository.GetURLSByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responseURLs []model.ShortenerURLSForUserResponse
	for k, v := range items {
		responseURLs = append(responseURLs, model.ShortenerURLSForUserResponse{
			OriginalURL: v,
			ShortURL:    s.generateResponseURL(k),
		})
	}

	return responseURLs, err
}

// DeleteUserURLS удаляет (помечает как удалённые) список ссылок, привязанных к пользователю.
func (s ShortenerService) DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error {
	if len(items) == 0 {
		return nil
	}

	return s.repository.DeleteUserURLS(ctx, items)
}

// ScheduleURLDeletion планирует отложенное удаление ссылок через канал.
func (s ShortenerService) ScheduleURLDeletion(ctx context.Context, items []model.URLToDelete) {
	go func() {
		for _, item := range items {
			s.deleteChan <- item
		}
	}()
}

// Run запускает фоновый процесс для пакетного удаления ссылок по таймеру или по лимиту.
func (s ShortenerService) Run(ctx context.Context, wg *sync.WaitGroup) {
	limit := 100
	batch := make([]model.URLToDelete, 0, limit)
	ticker := time.NewTicker(time.Second * 5)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer ticker.Stop()

		for {
			select {
			case item, ok := <-s.deleteChan:
				if !ok {
					if len(batch) > 0 {
						s.DeleteUserURLS(ctx, batch)
					}
					return
				}

				if len(batch) >= limit {
					s.DeleteUserURLS(ctx, batch)
					batch = batch[:0]
				}

				batch = append(batch, item)
			case <-ticker.C:
				s.DeleteUserURLS(ctx, batch)
				batch = batch[:0]
			case <-ctx.Done():
				if len(batch) > 0 {
					s.DeleteUserURLS(ctx, batch)
				}
				return
			}
		}
	}()
}

// Close - Закрывает канал
func (s ShortenerService) Close() {
	close(s.deleteChan)
}

// Stats — возвращает статистику по количеству URL и пользователей.
func (s ShortenerService) Stats(ctx context.Context) (model.StatsRespose, error) {
	return s.repository.Stats(ctx)
}
