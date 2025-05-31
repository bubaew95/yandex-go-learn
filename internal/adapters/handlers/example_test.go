package handlers

import (
	"bytes"
	"fmt"
	"github.com/bubaew95/yandex-go-learn/config"
	fileStorage "github.com/bubaew95/yandex-go-learn/internal/adapters/repository/filestorage"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"os"
)

func ExampleShortenerHandler_AddNewURL() {
	cfg := config.Config{
		FilePath: "data.json",
	}
	defer os.Remove(cfg.FilePath)

	shortenerDB, _ := storage.NewShortenerDB(cfg)
	shortenerRepository, err := fileStorage.NewShortenerRepository(*shortenerDB)
	if err != nil {
		fmt.Println("Ошибка")
		return
	}

	mockService := service.NewShortenerService(shortenerRepository, cfg)
	handler := NewShortenerHandler(mockService, cfg)

	route := chi.NewRouter()
	route.Post("/", handler.CreateURL)

	ts := httptest.NewServer(route)
	defer ts.Close()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	w := httptest.NewRecorder()

	handler.CreateURL(w, req)

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	// Output:
	// 201
}
