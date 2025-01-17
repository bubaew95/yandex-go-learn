package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/bubaew95/yandex-go-learn/internal/repository"
	"github.com/bubaew95/yandex-go-learn/internal/service"
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

	data := make(map[string]string)

	cfg := config.NewConfig()
	shortenerRepository := repository.NewShortenerRepository(data, cfg.BaseURL)
	shortenerService := service.NewShortenerService(shortenerRepository)
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
			url:  "/WzYAhpnS",
			want: want{
				contentType: "text/html",
				statusCode:  http.StatusTemporaryRedirect,
				location:    "https://practicum.yandex.ru/",
			},
		},

		{
			name: "Bad request test",
			url:  "/WzYAhS",
			want: want{
				contentType: "text/html",
				statusCode:  http.StatusBadRequest,
				location:    "https://practicum.yandex.ru/learn",
			},
		},
	}

	data := map[string]string{
		"WzYAhpnS": "https://practicum.yandex.ru/",
		"WzYAhS":   "https://practicum.yandex.ru/learn",
	}

	shortenerRepository := repository.NewShortenerRepository(data, "http://test.local")
	shortenerService := service.NewShortenerService(shortenerRepository)
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
