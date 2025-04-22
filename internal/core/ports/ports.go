// Package ports содержит интерфейсы (порты) для определения внешнего поведения
// хранилищ и сервисов в архитектуре приложения, следуя принципам Clean Architecture.
package ports

import (
	"context"

	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

// ShortenerRepository определяет контракт для репозитория сокращённых URL.
// Этот интерфейс реализуется различными адаптерами хранилищ (например, файловая система, PostgreSQL).
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
}

// ShortenerService определяет бизнес-логику сервиса сокращения ссылок.
// Включает в себя генерацию ссылок, работу с пользователями и отложенное удаление.
type ShortenerService interface {
	// GenerateURL генерирует короткий URL на основе оригинального.
	GenerateURL(ctx context.Context, url string, randomStringLength int) (string, error)

	// GetURLByID возвращает оригинальный URL по его сокращённому ID.
	GetURLByID(ctx context.Context, id string) (string, error)

	// GetURLByOriginalURL возвращает ID, соответствующий оригинальному URL.
	GetURLByOriginalURL(ctx context.Context, originalURL string) (string, bool)

	// InsertURLs добавляет множество URL и возвращает их короткие представления.
	InsertURLs(ctx context.Context, urls []model.ShortenerURLMapping) ([]model.ShortenerURLResponse, error)

	// GetURLSByUserID возвращает список сокращённых URL, принадлежащих пользователю.
	GetURLSByUserID(ctx context.Context, userID string) ([]model.ShortenerURLSForUserResponse, error)

	// DeleteUserURLS помечает ссылки как удалённые.
	DeleteUserURLS(ctx context.Context, items []model.URLToDelete) error

	// ScheduleURLDeletion планирует асинхронное удаление ссылок (например, через очередь).
	ScheduleURLDeletion(ctx context.Context, items []model.URLToDelete)

	// RandStringBytes генерирует случайную строку заданной длины (обычно для ID короткой ссылки).
	RandStringBytes(n int) string

	// Ping проверяет доступность сервиса (например, для liveness-проб).
	Ping(ctx context.Context) error
}
