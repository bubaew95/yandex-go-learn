package filestorage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

func BenchmarkShortenerRepository_InsertURLs(b *testing.B) {
	cfg := &config.Config{
		ServerAddress: "9090",
		BaseURL:       "http://test.local",
		FilePath:      "data.json",
	}

	defer os.Remove(cfg.FilePath)

	shortenerDB, _ := storage.NewShortenerDB(*cfg)

	shortener, _ := NewShortenerRepository(*shortenerDB)

	items := []model.ShortenerURLMapping{
		{
			CorrelationID: "rasf1D",
			OriginalURL:   "http://test.local",
		},
		{
			CorrelationID: "rasf2D",
			OriginalURL:   "http://test.local",
		},
	}

	b.ResetTimer()

	b.Run("insert", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			shortener.InsertURLs(context.Background(), items)
		}
	})

	b.Run("insert_2", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			shortener.InsertURLTwo(context.Background(), items)
		}
	})
}

func createTempStorageFile(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	return filepath.Join(dir, "shortener_test.json")
}

func TestShortenerRepository_SetAndGet(t *testing.T) {
	file := createTempStorageFile(t)

	cfg := config.Config{FilePath: file}
	db, err := storage.NewShortenerDB(cfg)
	require.NoError(t, err)

	repo, err := NewShortenerRepository(*db)
	require.NoError(t, err)

	err = repo.SetURL(context.Background(), "abc123", "https://example.com")
	require.NoError(t, err)

	// Get by ID
	got, err := repo.GetURLByID(context.Background(), "abc123")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", got)

	// Get by OriginalURL
	id, found := repo.GetURLByOriginalURL(context.Background(), "example.com")
	assert.True(t, found)
	assert.Equal(t, "abc123", id)
}

func TestShortenerRepository_InsertURLs(t *testing.T) {
	file := createTempStorageFile(t)

	cfg := config.Config{FilePath: file}
	db, err := storage.NewShortenerDB(cfg)
	require.NoError(t, err)

	repo, err := NewShortenerRepository(*db)
	require.NoError(t, err)

	urls := []model.ShortenerURLMapping{
		{CorrelationID: "id1", OriginalURL: "https://a.com"},
		{CorrelationID: "id2", OriginalURL: "https://b.com"},
	}

	err = repo.InsertURLs(context.Background(), urls)
	require.NoError(t, err)

	got1, err := repo.GetURLByID(context.Background(), "id1")
	require.NoError(t, err)
	assert.Equal(t, "https://a.com", got1)

	got2, err := repo.GetURLByID(context.Background(), "id2")
	require.NoError(t, err)
	assert.Equal(t, "https://b.com", got2)
}

func TestShortenerRepository_Ping(t *testing.T) {
	file := createTempStorageFile(t)

	cfg := config.Config{FilePath: file}
	db, err := storage.NewShortenerDB(cfg)
	require.NoError(t, err)

	repo, err := NewShortenerRepository(*db)
	require.NoError(t, err)

	err = repo.Ping(context.Background())
	require.NoError(t, err)
}

func TestShortenerRepository_InsertURLTwo(t *testing.T) {
	file := createTempStorageFile(t)

	cfg := config.Config{FilePath: file}
	db, err := storage.NewShortenerDB(cfg)
	require.NoError(t, err)

	repo, err := NewShortenerRepository(*db)
	require.NoError(t, err)

	// Сначала сохраним
	err = repo.SetURL(context.Background(), "id1", "https://old.com")
	require.NoError(t, err)

	// Попробуем обновить
	update := []model.ShortenerURLMapping{
		{CorrelationID: "id1", OriginalURL: "https://new.com"},
		{CorrelationID: "missing", OriginalURL: "https://should-not-be-added.com"},
	}

	err = repo.InsertURLTwo(context.Background(), update)
	require.NoError(t, err)

	// Должен быть обновлён
	got, err := repo.GetURLByID(context.Background(), "id1")
	require.NoError(t, err)
	assert.Equal(t, "https://new.com", got)

	// Не должен быть добавлен
	_, err = repo.GetURLByID(context.Background(), "missing")
	assert.Error(t, err)
}

func TestShortenerRepository_NotImplemented(t *testing.T) {
	file := createTempStorageFile(t)

	cfg := config.Config{FilePath: file}
	db, err := storage.NewShortenerDB(cfg)
	require.NoError(t, err)

	repo, err := NewShortenerRepository(*db)
	require.NoError(t, err)

	// DeleteUserURLS
	err = repo.DeleteUserURLS(context.Background(), nil)
	require.NoError(t, err)

	// GetURLSByUserID
	urls, err := repo.GetURLSByUserID(context.Background(), "any")
	require.NoError(t, err)
	assert.Nil(t, urls)
}

func TestShortenerRepository_Stats(t *testing.T) {
	file := createTempStorageFile(t)
	cfg := config.Config{FilePath: file}
	db, err := storage.NewShortenerDB(cfg)
	require.NoError(t, err)

	repo, err := NewShortenerRepository(*db)
	require.NoError(t, err)

	stats, err := repo.Stats(context.Background())
	require.NoError(t, err)

	assert.Equal(t, 0, stats.Users)
}
