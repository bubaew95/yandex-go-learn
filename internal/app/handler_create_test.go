package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bubaew95/yandex-go-learn/config"
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
	app := NewApp(cfg)
	app.Routers()

	ts := httptest.NewServer(&app.Router)
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
