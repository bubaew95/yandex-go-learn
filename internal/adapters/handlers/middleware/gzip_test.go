package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGZipMiddleware(t *testing.T) {
	t.Run("decompress gzip request and compress gzip response", func(t *testing.T) {
		originalBody := `{"message":"hello world"}`
		var compressedBody bytes.Buffer

		gz := gzip.NewWriter(&compressedBody)
		_, err := gz.Write([]byte(originalBody))
		require.NoError(t, err)
		require.NoError(t, gz.Close())

		handler := GZipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, originalBody, string(body))

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"ok":true}`))
		}))

		req := httptest.NewRequest(http.MethodPost, "/", &compressedBody)
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		gr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		defer gr.Close()

		decompressed, err := io.ReadAll(gr)
		require.NoError(t, err)
		assert.JSONEq(t, `{"ok":true}`, string(decompressed))
	})

	t.Run("no gzip headers — plain pass-through", func(t *testing.T) {
		handler := GZipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Equal(t, "plain text", string(body))

			w.Write([]byte("plain response"))
		}))

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("plain text"))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Empty(t, resp.Header.Get("Content-Encoding"))

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "plain response", string(respBody))
	})
}
