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

var (
	rsaPrivateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	rsaPublicKey     = &rsaPrivateKey.PublicKey
)

// New algorithms
func EncodeUserIDRSA(userID string) (string, error) {
	label := []byte("") // Optional label
	hash := sha256.New()

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, rsaPublicKey, []byte(userID), label)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecodeUserIDRSA(userId string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(userId)
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

func ValidateUserID(cookie *http.Cookie) bool {
	userID := cookie.Value
	signUserID, _ := DecodeUserID(userID)

	return userID != signUserID
}

func GenerateUserID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
