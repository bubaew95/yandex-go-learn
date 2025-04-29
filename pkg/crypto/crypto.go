// Package crypto предоставляет утилиты для генерации, кодирования и декодирования идентификаторов пользователей.
// Поддерживаются симметричное (AES) и асимметричное (RSA) шифрование.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type ctxKey string

// KeyUserID — ключ, используемый для хранения/извлечения ID пользователя из контекста.
const KeyUserID ctxKey = "user_id"

var (
	secretKey = "x35k9f"
)

// DecodeUserID расшифровывает закодированный идентификатор пользователя, зашифрованный с использованием AES.
// Возвращает оригинальное строковое значение.
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

// EncodeUserID шифрует идентификатор пользователя с использованием AES и возвращает hex-представление.
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

var (
	rsaPrivateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	rsaPublicKey     = &rsaPrivateKey.PublicKey
)

// EncodeUserIDRSA шифрует идентификатор пользователя с помощью RSA-OAEP.
// Возвращает base64-представление зашифрованных данных.
func EncodeUserIDRSA(userID string) (string, error) {
	label := []byte("") // Optional label
	hash := sha256.New()

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, rsaPublicKey, []byte(userID), label)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecodeUserIDRSA расшифровывает идентификатор пользователя, зашифрованный с помощью RSA-OAEP.
func DecodeUserIDRSA(userID string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(userID)
	if err != nil {
		return "", err
	}

	label := []byte("")
	hash := sha256.New()

	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, rsaPrivateKey, ciphertext, label)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// IsInvalidUserID проверяет соответствие значения cookie и декодированного ID.
// Возвращает true, если они не совпадают (что может указывать на подделку).
func IsInvalidUserID(cookie *http.Cookie) bool {
	userID := cookie.Value
	signUserID, _ := DecodeUserID(userID)

	return userID != signUserID
}

// GenerateUserID генерирует уникальный ID пользователя на основе текущего времени.
func GenerateUserID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
