package handlers

import (
	"bytes"
	"fmt"
	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/repository/postgres"
	"github.com/bubaew95/yandex-go-learn/internal/core/service"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
)

func ExampleShortenerHandler_AddNewURL() {
	cfg := config.Config{}

	shortenerRepository, err := postgres.NewShortenerRepository(cfg)
	if err != nil {
		logger.Log.Fatal("Database initialization error", zap.Error(err))
	}

	shortenerService := service.NewShortenerService(shortenerRepository, cfg)

	handler := NewShortenerHandler(shortenerService)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	w := httptest.NewRecorder()

	handler.CreateURL(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	fmt.Println(res.StatusCode)
	fmt.Println(string(body))

	// Output:
	// 201
	// http://short.url/abc123
}
