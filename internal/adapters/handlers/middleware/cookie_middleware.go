package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"go.uber.org/zap"
)

type ctxKey string

const KeyUserID ctxKey = "user_id"

var (
	secretKey = []byte("testss")
)

type favContextKey string

func CookieMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			cookieValue string
		)

		cookieUserID, err := r.Cookie("user_id")
		if err != nil || cookieUserID.Value == "" {
			userID := generateUserID()
			cookieValue = signUserID(userID)
		} else {
			isValid := validateUserID(cookieUserID)
			if !isValid {
				userID := generateUserID()
				cookieValue = signUserID(userID)
			} else {
				cookieValue = cookieUserID.Value
			}
		}

		logger.Log.Debug("user ud", zap.String("user_id", cookieValue))
		ctx := context.WithValue(r.Context(), KeyUserID, cookieValue)
		nRequest := r.WithContext(ctx)

		cookie := &http.Cookie{
			Name:  "user_id",
			Value: cookieValue,
		}

		// http.SetCookie(w, cookie)
		nRequest.AddCookie(cookie)
		h.ServeHTTP(w, nRequest)
	})
}

func signUserID(userID string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(userID))

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func validateUserID(cookie *http.Cookie) bool {
	userID := cookie.Value
	signUserID := signUserID(userID)

	return userID != signUserID
}

func generateUserID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
