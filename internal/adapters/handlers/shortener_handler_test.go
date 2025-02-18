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
	fileStorage "github.com/bubaew95/yandex-go-learn/internal/adapters/repository/filestorage"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/repository/postgres/mock"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/internal/core/service"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerCreate(t *testing.T) {
	t.Parallel()

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
	shortenerRepository, err := fileStorage.NewShortenerRepository(*shortenerDB)
	require.NoError(t, err)

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
	t.Parallel()

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

	shortenerDB.Save(&model.ShortenURL{
		UUID:        1,
		ShortURL:    "WzYAhS",
		OriginalURL: "https://practicum.yandex.ru/learn",
	})
	shortenerDB.Save(&model.ShortenURL{
		UUID:        2,
		ShortURL:    "WzYAhSs",
		OriginalURL: "https://practicum.yandex.ru/learn",
	})

	shortenerRepository, err := fileStorage.NewShortenerRepository(*shortenerDB)
	require.NoError(t, err)

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
	t.Parallel()

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

	shortenerRepository, err := fileStorage.NewShortenerRepository(*shortenerDB)
	require.NoError(t, err)

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

			var responseModel model.ShortenerResponse
			err = json.Unmarshal(respBody, &responseModel)
			require.NoError(t, err)

			assert.NotEmpty(t, responseModel.Result)
			assert.Equal(t, resp.StatusCode, tt.want.status)
			assert.Equal(t, resp.Header.Get("content-length"), tt.want.contentLength)
			assert.Equal(t, resp.Header.Get("content-type"), tt.want.contentType)
		})
	}
}

func TestHandlerBatch(t *testing.T) {
	t.Parallel()

	type want struct {
		status int
		result string
	}

	tests := []struct {
		name string
		data string
		want want
	}{
		{
			name: "Add urls success",
			data: `[ { "correlation_id": "test-1", "original_url": "http://google.com" }, { "correlation_id": "test-2", "original_url": "http://yandex.ru" }, { "correlation_id": "test-3", "original_url": "http://yandex.ru" } ]`,
			want: want{
				status: http.StatusCreated,
				result: `[ { "correlation_id": "test-1", "short_url": "https://site.local/test-1" }, { "correlation_id": "test-2", "short_url": "https://site.local/test-2" }, { "correlation_id": "test-3", "short_url": "https://site.local/test-3" } ]`,
			},
		},
		{
			name: "Dublicate CorrelationId",
			data: `[{ "correlation_id": "test-1", "original_url": "http://google.com" }, { "correlation_id": "test-1", "original_url": "http://yandex.ru" }]`,
			want: want{
				status: http.StatusCreated,
				result: `[{ "correlation_id": "test-1", "short_url": "https://site.local/test-1" }, { "correlation_id": "test-1", "short_url": "https://site.local/test-1" }]`,
			},
		},
	}

	cfg := &config.Config{
		BaseURL: "https://site.local",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shortenerRepository := mock.NewMockShortenerRepository(ctrl)
	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)
	shortenerHandler := NewShortenerHandler(shortenerService)

	router := chi.NewRouter()
	router.Post("/api/shorten/batch", shortenerHandler.Batch)

	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var items []model.ShortenerURLMapping
			err := json.Unmarshal([]byte(tt.data), &items)
			require.NoError(t, err)

			shortenerRepository.EXPECT().
				InsertURLs(gomock.Any(), items).
				Return(nil)

			req, err := http.Post(ts.URL+"/api/shorten/batch", "application/json", strings.NewReader(tt.data))
			require.NoError(t, err)
			defer req.Body.Close()

			respBody, err := io.ReadAll(req.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.status, req.StatusCode)
			if tt.want.result != "" {
				assert.JSONEq(t, tt.want.result, string(respBody))
			}
		})
	}
}

// func TestUserURLS(t *testing.T) {
// 	t.Parallel()

// 	type want struct {
// 		status int
// 		result string
// 	}

// 	tests := []struct {
// 		name   string
// 		path   string
// 		data   string
// 		method string
// 		userID string
// 		want   want
// 	}{
// 		{
// 			name:   "Add url for user",
// 			path:   `/`,
// 			method: http.MethodPost,
// 			data:   `http://google.com`,
// 			want: want{
// 				status: http.StatusCreated,
// 			},
// 			userID: "user_id",
// 		},
// 		{
// 			name:   "Get User urls",
// 			path:   `/api/user/urls`,
// 			method: http.MethodGet,
// 			want: want{
// 				status: http.StatusCreated,
// 				result: ``,
// 			},
// 		},
// 	}

// 	cfg := &config.Config{
// 		BaseURL:     "https://site.local",
// 		DataBaseDSN: "host=127.0.0.1 user=admin password=admin dbname=yandex sslmode=disable",
// 	}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// shortenerRepository := mock.NewMockShortenerRepository(ctrl)
// 	shortenerRepository, err := postgres.NewShortenerRepository(*cfg)
// 	require.NoError(t, err)

// 	shortenerService := service.NewShortenerService(shortenerRepository, *cfg)
// 	shortenerHandler := NewShortenerHandler(shortenerService)

// 	router := chi.NewRouter()
// 	router.Post("/", shortenerHandler.CreateURL)
// 	router.Get("/api/user/urls", shortenerHandler.GetUserURLS)

// 	ts := httptest.NewServer(router)
// 	defer ts.Close()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req, err := http.NewRequest(tt.method, ts.URL+tt.path, strings.NewReader(tt.data))
// 			require.NoError(t, err)
// 			defer req.Body.Close()

// 			resp, err := ts.Client().Do(req)
// 			require.NoError(t, err)
// 			defer resp.Body.Close()

// 			fmt.Println("user_id", req.Context().Value(crypto.KeyUserID))

// 			respBody, err := io.ReadAll(resp.Body)
// 			require.NoError(t, err)

// 			fmt.Println(tt.path, string(respBody))
// 		})
// 	}
// }
