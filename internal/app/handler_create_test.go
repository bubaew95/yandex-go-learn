package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerCreate(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		method      string
		url         string
	}

	tests := []struct {
		name string
		data string
		want want
	}{
		{
			name: "Simple url",
			data: "https://practicum.yandex.ru/",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				method:      http.MethodPost,
				url:         "http://localhost:8080/WzYAhpnS",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()
			w.Header().Set("content-type", "text/plain")

			urls := make(map[string]string)
			handler := http.HandlerFunc(CreateURL(urls))
			handler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("content-type"))

			data, err := io.ReadAll(result.Body)
			require.NoError(t, err, "Ошибка получения данных")

			err = result.Body.Close()
			require.NoError(t, err, "Ошибка закрытия подключения")

			assert.NotEmpty(t, data)
		})
	}

}
