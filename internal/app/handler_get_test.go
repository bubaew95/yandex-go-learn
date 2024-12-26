package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerGet(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		method      string
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
			data: map[string]string{
				"WzYAhpnS": "https://practicum.yandex.ru/",
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusTemporaryRedirect,
				method:      http.MethodGet,
				location:    "https://practicum.yandex.ru/",
			},
		},

		{
			name: "Bad request test",
			url:  "/WzYAhS",
			data: map[string]string{
				"WzYAhpnS": "https://practicum.yandex.ru/",
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
				method:      http.MethodGet,
				location:    "https://practicum.yandex.ru/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.want.method, tt.url, nil)
			w := httptest.NewRecorder()
			w.Header().Set("content-type", tt.want.contentType)
			w.Header().Set("location", tt.want.location)
			w.WriteHeader(tt.want.statusCode)

			handler := http.HandlerFunc(GetURL(tt.data))
			handler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("content-type"))
			assert.Equal(t, tt.want.location, result.Header.Get("location"))
		})
	}
}
