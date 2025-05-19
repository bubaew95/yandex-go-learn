package handlers

import (
	"encoding/json"
	"errors"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/constants"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bubaew95/yandex-go-learn/config"
	fileStorage "github.com/bubaew95/yandex-go-learn/internal/adapters/repository/filestorage"
	"github.com/bubaew95/yandex-go-learn/internal/adapters/storage"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
	"github.com/bubaew95/yandex-go-learn/internal/core/service"
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
				contentType: "text/html; charset=utf-8",
				statusCode:  http.StatusTemporaryRedirect,
				location:    "https://practicum.yandex.ru/",
			},
		},
		{
			name: "Bad request test",
			url:  "/WzYAhSs",
			want: want{
				contentType: "text/html; charset=utf-8",
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

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.NotEmpty(t, respBody)
			assert.NotEmpty(t, resp.Header.Get("content-type"))
		})
	}
}

func TestHandlerAddNewURLFromJson(t *testing.T) {
	t.Parallel()

	type want struct {
		contentType string
		status      int
	}

	tests := []struct {
		name       string
		request    string
		mockErr    error
		mockResult string
		originURL  string
		want       want
	}{
		{
			name:       "Add new URL successfully",
			request:    `{"url": "https://practicum.yandex.ru"}`,
			mockResult: "http://short.url/abc123",
			want: want{
				contentType: "application/json",
				status:      http.StatusCreated,
			},
		},
		{
			name:      "Conflict - URL already exists",
			request:   `{"url": "https://practicum.yandex.ru"}`,
			mockErr:   constants.ErrUniqueIndex,
			originURL: "http://short.url/existing",
			want: want{
				contentType: "application/json",
				status:      http.StatusConflict,
			},
		},
		{
			name:    "Invalid JSON",
			request: `invalid json`,
			want: want{
				status: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortenerService := NewMockShortenerService(t)
			handler := NewShortenerHandler(shortenerService)

			router := chi.NewRouter()
			router.Post("/api/shorten", handler.AddNewURL)

			ts := httptest.NewServer(router)
			defer ts.Close()

			if tt.mockErr != nil {
				shortenerService.On("GenerateURL", mock.Anything, "https://practicum.yandex.ru", mock.Anything).
					Return("", tt.mockErr).
					Once()

				if tt.originURL != "" {
					shortenerService.On("GetURLByOriginalURL", mock.Anything, "https://practicum.yandex.ru").
						Return(tt.originURL, true).
						Once()
				}
			} else if tt.mockResult != "" {
				shortenerService.On("GenerateURL", mock.Anything, "https://practicum.yandex.ru", mock.Anything).
					Return(tt.mockResult, nil).
					Once()
			}

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", strings.NewReader(tt.request))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.status, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if tt.want.status == http.StatusCreated || tt.want.status == http.StatusConflict {
				var response model.ShortenerResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				expected := tt.mockResult
				if tt.mockErr != nil {
					expected = tt.originURL
				}
				assert.Equal(t, expected, response.Result)
			}
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
		name    string
		data    string
		want    want
		isError bool
	}{
		{
			name: "Add urls success",
			data: `[ { "correlation_id": "test-1", "original_url": "http://google.com" }, { "correlation_id": "test-2", "original_url": "http://yandex.ru" }, { "correlation_id": "test-3", "original_url": "http://yandex.ru" } ]`,
			want: want{
				status: http.StatusCreated,
				result: `[ { "correlation_id": "test-1", "short_url": "https://site.local/test-1" }, { "correlation_id": "test-2", "short_url": "https://site.local/test-2" }, { "correlation_id": "test-3", "short_url": "https://site.local/test-3" } ]`,
			},
			isError: false,
		},
		{
			name: "Dublicate CorrelationId",
			data: `[{ "correlation_id": "test-1", "original_url": "http://google.com" }, { "correlation_id": "test-1", "original_url": "http://yandex.ru" }]`,
			want: want{
				status: http.StatusCreated,
				result: `[{ "correlation_id": "test-1", "short_url": "https://site.local/test-1" }, { "correlation_id": "test-1", "short_url": "https://site.local/test-1" }]`,
			},
			isError: false,
		},
		{
			name: "Invalidate json",
			data: `[{ "test-1", "original_url": "http://google.com" }]`,
			want: want{
				status: http.StatusInternalServerError,
				result: ``,
			},
			isError: true,
		},
	}

	shortenerService := NewMockShortenerService(t)
	shortenerHandler := NewShortenerHandler(shortenerService)

	router := chi.NewRouter()
	router.Post("/api/shorten/batch", shortenerHandler.Batch)

	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.isError {
				var items []model.ShortenerURLMapping
				err := json.Unmarshal([]byte(tt.data), &items)
				require.NoError(t, err)

				var respItems []model.ShortenerURLResponse
				for _, item := range items {
					respItems = append(respItems, model.ShortenerURLResponse{
						CorrelationID: item.CorrelationID,
						ShortURL:      "https://site.local/" + item.CorrelationID,
					})
				}

				shortenerService.On("InsertURLs", mock.Anything, items).
					Return(respItems, nil)
			}

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

func TestShortenerHandler_GetUserURLS(t *testing.T) {
	t.Parallel()

	type want struct {
		status      int
		contentType string
	}

	tests := []struct {
		name       string
		cookie     *http.Cookie
		mockUserID string
		mockURLs   []model.ShortenerURLSForUserResponse
		mockErr    error
		want       want
	}{
		{
			name:   "No cookie present",
			cookie: nil,
			want: want{
				status: http.StatusNoContent,
			},
		},
		{
			name:       "Service returns error",
			cookie:     &http.Cookie{Name: "user_id", Value: "user123"},
			mockUserID: "user123",
			mockErr:    errors.New("service failure"),
			want: want{
				status: http.StatusInternalServerError,
			},
		},
		{
			name:       "No URLs found",
			cookie:     &http.Cookie{Name: "user_id", Value: "user123"},
			mockUserID: "user123",
			mockURLs:   nil,
			want: want{
				status: http.StatusNoContent,
			},
		},
		{
			name:       "URLs found",
			cookie:     &http.Cookie{Name: "user_id", Value: "user123"},
			mockUserID: "user123",
			mockURLs: []model.ShortenerURLSForUserResponse{
				{ShortURL: "http://short.url/abc", OriginalURL: "http://example.com"},
				{ShortURL: "http://short.url/def", OriginalURL: "http://example.org"},
			},
			want: want{
				status:      http.StatusOK,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		tt := tt // захват переменной
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := NewMockShortenerService(t)
			handler := ShortenerHandler{service: mockService}

			router := chi.NewRouter()
			router.Get("/api/user/urls", handler.GetUserURLS)
			ts := httptest.NewServer(router)
			defer ts.Close()

			if tt.cookie != nil {
				mockService.On("GetURLSByUserID", mock.Anything, tt.mockUserID).
					Return(tt.mockURLs, tt.mockErr).
					Once()
			}

			req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/user/urls", nil)
			require.NoError(t, err)

			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.status, resp.StatusCode)
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			}

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if tt.want.status == http.StatusOK {
				var urls []model.ShortenerURLSForUserResponse
				err = json.Unmarshal(body, &urls)
				require.NoError(t, err)
				assert.Equal(t, tt.mockURLs, urls)
			} else {
				assert.Empty(t, body)
			}
		})
	}
}

func TestShortenerHandler_DeleteUserURLS(t *testing.T) {
	t.Parallel()

	type want struct {
		status int
	}

	tests := []struct {
		name       string
		cookie     *http.Cookie
		body       string
		mockCalled bool // хотим проверить, что ScheduleURLDeletion вызвался
		want       want
	}{
		{
			name:       "No cookie returns 204",
			cookie:     nil,
			body:       `["http://short.url/abc"]`,
			mockCalled: false,
			want: want{
				status: http.StatusNoContent,
			},
		},
		{
			name:       "Invalid JSON returns 500",
			cookie:     &http.Cookie{Name: "user_id", Value: "user123"},
			body:       `invalid json`,
			mockCalled: false,
			want: want{
				status: http.StatusInternalServerError,
			},
		},
		{
			name:       "Valid request schedules deletion and returns 202",
			cookie:     &http.Cookie{Name: "user_id", Value: "user123"},
			body:       `["http://short.url/abc", "http://short.url/def"]`,
			mockCalled: true,
			want: want{
				status: http.StatusAccepted,
			},
		},
	}

	for _, tt := range tests {
		tt := tt // захват переменной
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := NewMockShortenerService(t)
			handler := ShortenerHandler{service: mockService}

			router := chi.NewRouter()
			router.Delete("/api/user/urls", handler.DeleteUserURLS)
			ts := httptest.NewServer(router)
			defer ts.Close()

			if tt.mockCalled {
				// Ожидаем вызов ScheduleURLDeletion с нужными параметрами
				mockService.On("ScheduleURLDeletion", mock.Anything, mock.MatchedBy(func(items []model.URLToDelete) bool {
					if len(items) != 2 {
						return false
					}
					return items[0].ShortLink == "http://short.url/abc" && items[0].UserID == "user123" &&
						items[1].ShortLink == "http://short.url/def" && items[1].UserID == "user123"
				})).Return().Once()
			}

			req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/user/urls", strings.NewReader(tt.body))
			require.NoError(t, err)

			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.status, resp.StatusCode)

			if tt.mockCalled {
				mockService.AssertExpectations(t)
			}
		})
	}
}
