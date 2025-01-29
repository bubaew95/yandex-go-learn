package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/models"
	"github.com/bubaew95/yandex-go-learn/internal/repository"
	"github.com/bubaew95/yandex-go-learn/internal/service"
	"github.com/bubaew95/yandex-go-learn/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerCreate(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		method      string
	}

	tests := []struct {
		name string
		data string
		want want
		err  bool
	}{
		{
			name: "Simple url",
			data: "https://practicum.yandex.ru/",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusCreated,
				method:      http.MethodPost,
			},
			err: false,
		},
		{
			name: "Data is empty",
			data: "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				method:      http.MethodPost,
			},

			err: true,
		},
	}

	cfg := config.NewConfig()

	shortenerDB, _ := storage.NewShortenerDB(*cfg)
	shortenerRepository := repository.NewShortenerRepository(*shortenerDB)
	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)
	shortenerHandler := NewShortenerHandler(shortenerService)

	route := chi.NewRouter()
	route.Post("/", shortenerHandler.CreateURL)

	ts := httptest.NewServer(route)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", ts.URL, strings.NewReader(tt.data))
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, resp.StatusCode, tt.want.statusCode)

			if false == tt.err {
				assert.Equal(t, resp.Header.Get("content-type"), tt.want.contentType)
				assert.Equal(t, respBody, respBody)
			}
		})
	}
}

func TestHandlerGet(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		location    string
	}

	tests := []struct {
		name string
		url  string
		data map[string]string
		want want
		err  bool
	}{
		{
			name: "Simple test",
			url:  "/WzYAhS",
			want: want{
				contentType: "text/html",
				statusCode:  http.StatusTemporaryRedirect,
				location:    "https://practicum.yandex.ru/",
			},
		},
		{
			name: "Bad request test",
			url:  "/WzYAhSs",
			want: want{
				contentType: "text/html",
				statusCode:  http.StatusBadRequest,
				location:    "https://practicum.yandex.ru/learn",
			},
		},
	}

	cfg := &config.Config{
		Port:     "9090",
		BaseURL:  "http://test.local",
		FilePath: "data.json",
	}

	shortenerDB, _ := storage.NewShortenerDB(*cfg)
	defer os.Remove(cfg.FilePath)

	shortenerDB.Save(&models.ShortenURL{
		UUID:        1,
		ShortURL:    "WzYAhS",
		OriginalURL: "https://practicum.yandex.ru/learn",
	})
	shortenerDB.Save(&models.ShortenURL{
		UUID:        2,
		ShortURL:    "WzYAhSs",
		OriginalURL: "https://practicum.yandex.ru/learn",
	})

	shortenerRepository := repository.NewShortenerRepository(*shortenerDB)
	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)
	shortenerHandler := NewShortenerHandler(shortenerService)

	route := chi.NewRouter()
	route.Get("/{id}", shortenerHandler.GetURL)

	ts := httptest.NewServer(route)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tt.url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, resp.Header.Get("content-type"), tt.want.contentType)
		})
	}
}

func TestHandlerAddNewURLFromJson(t *testing.T) {
	type want struct {
		contentType   string
		contentLength string
		method        string
		status        int
	}

	tests := []struct {
		name string
		data string
		want want
	}{
		{
			name: "Add new test from json",
			data: `{ "url": "https://practicum.yandex.ru"}`,
			want: want{
				contentLength: "40",
				contentType:   "application/json",
				method:        http.MethodPost,
				status:        http.StatusCreated,
			},
		},
	}

	cfg := &config.Config{
		Port:     "9090",
		BaseURL:  "http://test.local",
		FilePath: "data.json",
	}

	shortenerDB, _ := storage.NewShortenerDB(*cfg)
	defer os.Remove(cfg.FilePath)

	shortenerRepository := repository.NewShortenerRepository(*shortenerDB)
	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)
	shortenerHandler := NewShortenerHandler(shortenerService)

	router := chi.NewRouter()
	router.Post("/api/shorten", shortenerHandler.AddNewURL)

	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.want.method, ts.URL+"/api/shorten", bytes.NewBuffer([]byte(tt.data)))
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var responseModel models.ShortenerResponse
			err = json.Unmarshal(respBody, &responseModel)
			require.NoError(t, err)

			assert.NotEmpty(t, responseModel.Result)
			assert.Equal(t, resp.StatusCode, tt.want.status)
			assert.Equal(t, resp.Header.Get("content-length"), tt.want.contentLength)
			assert.Equal(t, resp.Header.Get("content-type"), tt.want.contentType)
		})
	}
}
