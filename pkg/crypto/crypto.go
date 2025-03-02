package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type ctxKey string

const KeyUserID ctxKey = "user_id"

var (
	secretKey = "x35k9f"
)

func DecodeUserID(userID string) (string, error) {
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

func EncodeUserID(userID string) (string, error) {
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

func ValidateUserID(cookie *http.Cookie) bool {
	userID := cookie.Value
	signUserID, _ := DecodeUserID(userID)

	return userID != signUserID
}

func GenerateUserID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
