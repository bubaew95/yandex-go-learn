package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bubaew95/yandex-go-learn/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	cfg := config.NewTestConfig([]string{"-a", ":8081", "-b", "http://test.local"})
	app := NewApp(cfg)
	app.Routers()
	fmt.Print(cfg)

	app.URLs = map[string]string{
		"WzYAhpnS": "https://practicum.yandex.ru/",
		"WzYAhS":   "https://practicum.yandex.ru/learn",
	}

	ts := httptest.NewServer(&app.Router)
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
