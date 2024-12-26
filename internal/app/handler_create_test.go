package app

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

			urls := make(map[string]string)
			h := http.HandlerFunc(CreateUrl(urls))
			h(w, request)

			result := w.Result()
			data, err := io.ReadAll(result.Body)
			require.NoError(t, err, "Ошибка получения данных")

			fmt.Println(data)
		})
	}

}
