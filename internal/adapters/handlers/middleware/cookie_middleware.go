package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/bubaew95/yandex-go-learn/internal/adapters/logger"
	"go.uber.org/zap"
)

type ctxKey string

const KeyUserID ctxKey = "user_id"

var (
	secretKey = "x35k9f"
)

func CookieMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			cookieValue string
			userID      string
		)

		cookieUserID, err := r.Cookie("user_id")
		if err != nil || cookieUserID.Value == "" || !validateUserID(cookieUserID) {
			userID = generateUserID()
			cookieValue, err = encodeUserID(userID)
			if err != nil {
				logger.Log.Debug("Error encode user id", zap.String("user_id", userID))
			}
		} else {
			cookieValue = cookieUserID.Value
		}

		ctx := context.WithValue(r.Context(), KeyUserID, cookieValue)
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

func decodeUserID(userID string) (string, error) {
	aesgcm, nonce, err := aesGcm()
	if err != nil {
		return "", err
	}

	encrypted, err := hex.DecodeString(userID)
	if err != nil {
		return "", err
	}

	decrypted, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func encodeUserID(userID string) (string, error) {
	aesgcm, nonce, err := aesGcm()
	if err != nil {
		return "", err
	}

	dst := aesgcm.Seal(nil, nonce, []byte(userID), nil)
	return fmt.Sprintf("%x", dst), nil
}

func aesGcm() (cipher.AEAD, []byte, error) {
	key := sha256.Sum256([]byte(secretKey))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, nil, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]
	return aesgcm, nonce, nil
}

func validateUserID(cookie *http.Cookie) bool {
	userID := cookie.Value
	signUserID, _ := decodeUserID(userID)

	return userID != signUserID
}

func generateUserID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
