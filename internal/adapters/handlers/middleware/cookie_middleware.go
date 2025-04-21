package middleware

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"github.com/bubaew95/yandex-go-learn/pkg/crypto"
)

// CookieMiddleware — middleware, обеспечивающий наличие user_id в куках пользователя.
func CookieMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			cookieValue string
			userID      string
		)

		cookieUserID, err := r.Cookie("user_id")
		if err != nil || cookieUserID.Value == "" || !crypto.ValidateUserID(cookieUserID) {
			userID = crypto.GenerateUserID()
			cookieValue, err = crypto.EncodeUserID(userID)
			if err != nil {
				logger.Log.Debug("Error encode user id", zap.String("user_id", userID))
			}
		} else {
			cookieValue = cookieUserID.Value
		}

		ctx := context.WithValue(r.Context(), crypto.KeyUserID, cookieValue)
		nRequest := r.WithContext(ctx)

		cookie := &http.Cookie{
			Name:     "user_id",
			Value:    cookieValue,
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
		h.ServeHTTP(w, nRequest)
	})
}
