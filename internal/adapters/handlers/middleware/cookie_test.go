package middleware

import (
	"github.com/bubaew95/yandex-go-learn/pkg/crypto"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCookieMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		requestCookie *http.Cookie
		expectReuse   bool
	}{
		{
			name:        "no cookie — generate new",
			expectReuse: false,
		},
		{
			name:          "valid cookie — reuse it",
			requestCookie: &http.Cookie{Name: "user_id", Value: mustEncoded("valid-user-id")},
			expectReuse:   true,
		},
		{
			name:          "invalid cookie — generate new",
			requestCookie: &http.Cookie{Name: "user_id", Value: "bad-cookie"},
			expectReuse:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cookieValueFromContext string

			// handler сохраняет user_id из контекста
			handler := CookieMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				val := r.Context().Value(crypto.KeyUserID)
				require.NotNil(t, val)
				cookieValueFromContext = val.(string)
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.requestCookie != nil {
				req.AddCookie(tt.requestCookie)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			resp := w.Result()

			cookies := resp.Cookies()
			var userIDCookie *http.Cookie
			for _, c := range cookies {
				if c.Name == "user_id" {
					userIDCookie = c
					break
				}
			}

			require.NotNil(t, userIDCookie, "user_id cookie must be set")

			if tt.expectReuse {
				require.Equal(t, tt.requestCookie.Value, userIDCookie.Value)
				require.Equal(t, tt.requestCookie.Value, cookieValueFromContext)
			} else {
				require.NotEmpty(t, userIDCookie.Value)
				require.Equal(t, userIDCookie.Value, cookieValueFromContext)
			}
		})
	}
}

func mustEncoded(uid string) string {
	val, err := crypto.EncodeUserID(uid)
	if err != nil {
		panic(err)
	}
	return val
}
