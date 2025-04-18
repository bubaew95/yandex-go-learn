package filestorage

import (
	"context"
	"os"
	"testing"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

func BenchmarkShortenerRepository_InsertURLs(b *testing.B) {
	cfg := &config.Config{
		Port:     "9090",
		BaseURL:  "http://test.local",
		FilePath: "data.json",
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
