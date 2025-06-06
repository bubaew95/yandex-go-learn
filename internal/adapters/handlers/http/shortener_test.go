package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/constants"

	"github.com/stretchr/testify/mock"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/core/model"
)

func TestHandlerCreate(t *testing.T) {
	t.Parallel()

	type want struct {
		statusCode  int
		contentType string
		expectBody  string
	}

	tests := []struct {
		name      string
		data      string
		mockSetup func(service *MockShortenerService)
		want      want
	}{
		{
			name: "Simple url - Created",
			data: "https://practicum.yandex.ru",
			mockSetup: func(service *MockShortenerService) {
				service.On("GenerateURL", mock.Anything, "https://practicum.yandex.ru", mock.Anything).
					Return("http://localhost:8080/abc123", nil).
					Once()
			},
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
				expectBody:  "http://localhost:8080/abc123",
			},
		},
		{
			name: "Empty body - Bad Request",
			data: "",
			mockSetup: func(service *MockShortenerService) {
				// Ничего не вызывается
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name: "Duplicate URL - Conflict",
			data: "https://duplicate.example.com",
			mockSetup: func(service *MockShortenerService) {
				service.On("GenerateURL", mock.Anything, "https://duplicate.example.com", mock.Anything).
					Return("", constants.ErrUniqueIndex).
					Once()

				service.On("GetURLByOriginalURL", mock.Anything, "https://duplicate.example.com").
					Return("http://localhost:8080/dupl123", true).
					Once()
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
				expectBody:  "http://localhost:8080/dupl123",
			},
		},
		{
			name: "Internal error - Generate fails",
			data: "https://fail.example.com",
			mockSetup: func(service *MockShortenerService) {
				service.On("GenerateURL", mock.Anything, "https://fail.example.com", mock.Anything).
					Return("", errors.New("internal error")).
					Once()
			},
			want: want{
				statusCode:  http.StatusInternalServerError,
				contentType: "",
			},
		},
	}

	cfg := config.NewConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortenerService := NewMockShortenerService(t)
			if tt.mockSetup != nil {
				tt.mockSetup(shortenerService)
			}

			handler := NewShortenerHandler(shortenerService, *cfg)

			r := chi.NewRouter()
			r.Post("/", handler.CreateURL)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodPost, ts.URL, strings.NewReader(tt.data))
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			if tt.want.expectBody != "" {
				assert.Equal(t, tt.want.expectBody, string(body))
			}

			shortenerService.AssertExpectations(t)
		})
	}
}

func TestHandlerGet(t *testing.T) {
	t.Parallel()

	type want struct {
		statusCode int
		location   string
	}

	tests := []struct {
		name       string
		requestURL string
		mockSetup  func(service *MockShortenerService)
		want       want
	}{
		{
			name:       "Found URL - 302 Redirect",
			requestURL: "/abc123",
			mockSetup: func(service *MockShortenerService) {
				service.On("GetURLByID", mock.Anything, "abc123").
					Return("https://practicum.yandex.ru/", nil).
					Once()
			},
			want: want{
				statusCode: http.StatusOK,
				location:   "",
			},
		},
		{
			name:       "Deleted URL - 410 Gone",
			requestURL: "/deleted123",
			mockSetup: func(service *MockShortenerService) {
				service.On("GetURLByID", mock.Anything, "deleted123").
					Return("", constants.ErrIsDeleted).
					Once()
			},
			want: want{
				statusCode: http.StatusGone,
			},
		},
		{
			name:       "Not found - 404 Not Found",
			requestURL: "/notfound",
			mockSetup: func(service *MockShortenerService) {
				service.On("GetURLByID", mock.Anything, "notfound").
					Return("", nil).
					Once()
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:       "Error from service - 404 Not Found",
			requestURL: "/errorcase",
			mockSetup: func(service *MockShortenerService) {
				service.On("GetURLByID", mock.Anything, "errorcase").
					Return("", errors.New("some error")).
					Once()
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}

	cfg := &config.Config{
		ServerAddress: "9090",
		BaseURL:       "http://test.local",
		FilePath:      "mock.json",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockShortenerService(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}

			handler := NewShortenerHandler(mockService, *cfg)

			router := chi.NewRouter()
			router.Get("/{id}", handler.GetURL)

			ts := httptest.NewServer(router)
			defer ts.Close()

			resp, err := http.Get(ts.URL + tt.requestURL)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			if tt.want.statusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
			} else {
				assert.Empty(t, resp.Header.Get("Location"))
			}

			mockService.AssertExpectations(t)
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

	cfg := &config.Config{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortenerService := NewMockShortenerService(t)
			handler := NewShortenerHandler(shortenerService, *cfg)

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

	cfg := &config.Config{}
	shortenerService := NewMockShortenerService(t)
	shortenerHandler := NewShortenerHandler(shortenerService, *cfg)

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

	cfg := &config.Config{}

	for _, tt := range tests {
		tt := tt // захват переменной
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := NewMockShortenerService(t)
			handler := ShortenerHandler{service: mockService, config: cfg}

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

	cfg := &config.Config{}

	for _, tt := range tests {
		tt := tt // захват переменной
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := NewMockShortenerService(t)
			handler := ShortenerHandler{service: mockService, config: cfg}

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

func TestHandlerStats(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		trustedSubnet string
		ipHeader      string
		expectedCode  int
		setupMock     func(service *MockShortenerService)
		expectBody    *model.StatsRespose
	}

	tests := []testCase{
		{
			name:          "OK — allowed IP in subnet",
			trustedSubnet: "127.0.0.1/32",
			ipHeader:      "127.0.0.1",
			expectedCode:  http.StatusOK,
			setupMock: func(service *MockShortenerService) {
				service.On("Stats", mock.Anything).
					Return(model.StatsRespose{Users: 5, URLs: 10}, nil).Once()
			},
			expectBody: &model.StatsRespose{Users: 5, URLs: 10},
		},
		{
			name:          "Forbidden — empty TrustedSubnet",
			trustedSubnet: "",
			ipHeader:      "127.0.0.1",
			expectedCode:  http.StatusForbidden,
		},
		{
			name:          "Internal error — invalid CIDR",
			trustedSubnet: "invalid-cidr",
			ipHeader:      "127.0.0.1",
			expectedCode:  http.StatusInternalServerError,
		},
		{
			name:          "Forbidden — IP not in subnet",
			trustedSubnet: "127.0.0.1/32",
			ipHeader:      "192.168.0.1",
			expectedCode:  http.StatusForbidden,
		},
		{
			name:          "Forbidden — missing IP header",
			trustedSubnet: "127.0.0.1/32",
			ipHeader:      "",
			expectedCode:  http.StatusForbidden,
		},
		{
			name:          "Internal error — Stats() fails",
			trustedSubnet: "127.0.0.1/32",
			ipHeader:      "127.0.0.1",
			expectedCode:  http.StatusInternalServerError,
			setupMock: func(service *MockShortenerService) {
				service.On("Stats", mock.Anything).
					Return(model.StatsRespose{}, errors.New("db error")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				TrustedSubnet: tt.trustedSubnet,
			}

			mockService := NewMockShortenerService(t)

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			handler := NewShortenerHandler(mockService, *cfg)

			router := chi.NewRouter()
			router.Get("/api/internal/stats", handler.Stats)

			ts := httptest.NewServer(router)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/internal/stats", nil)
			require.NoError(t, err)
			if tt.ipHeader != "" {
				req.Header.Set("X-Real-IP", tt.ipHeader)
			}

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.expectedCode == http.StatusOK && tt.expectBody != nil {
				var result model.StatsRespose
				err := json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				assert.Equal(t, *tt.expectBody, result)
			} else {
				body, _ := io.ReadAll(resp.Body)
				assert.Empty(t, string(body)) // тело пустое при ошибках
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandlerPing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mockError error
		wantCode  int
	}{
		{
			name:      "Ping success",
			mockError: nil,
			wantCode:  http.StatusOK,
		},
		{
			name:      "Ping failure",
			mockError: errors.New("db unavailable"),
			wantCode:  http.StatusInternalServerError,
		},
	}

	cfg := &config.Config{
		ServerAddress: "9090",
		BaseURL:       "http://test.local",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockShortenerService(t)
			mockService.On("Ping", mock.Anything).Return(tt.mockError).Once()

			handler := NewShortenerHandler(mockService, *cfg)

			r := chi.NewRouter()
			r.Get("/ping", handler.Ping)

			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, err := http.Get(ts.URL + "/ping")
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantCode, resp.StatusCode)

			mockService.AssertExpectations(t)
		})
	}
}
