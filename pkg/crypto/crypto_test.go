package crypto

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrypto(t *testing.T) {
	token := GenerateUserID()

	encodeUserID, err := EncodeUserID(token)
	require.NoError(t, err)

	decodeUserID, err := DecodeUserID(encodeUserID)
	require.NoError(t, err)

	encIDRSA, err := EncodeUserIDRSA(token)
	require.NoError(t, err)

	decodeIDRSA, err := DecodeUserIDRSA(encIDRSA)
	require.NoError(t, err)

	assert.NotEmpty(t, token)
	assert.True(t, decodeUserID == token)
	assert.True(t, decodeIDRSA == token)
}

func BenchmarkDecodeUserID(b *testing.B) {
	eUserID, _ := EncodeUserID("test")
	eIUserID, _ := EncodeUserIDRSA("test")

	b.Run("simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DecodeUserID(eUserID)
		}
	})

	b.Run("shuffle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DecodeUserIDRSA(eIUserID)
		}
	})
}

func BenchmarkEncodeUserID(b *testing.B) {
	b.Run("simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			EncodeUserID("test")
		}
	})

	b.Run("shuffle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			EncodeUserIDRSA("test")
		}
	})
}

func TestIsInvalidUserID(t *testing.T) {
	t.Run("Valid signed cookie", func(t *testing.T) {
		userID := GenerateUserID()

		encoded, err := EncodeUserID(userID)
		require.NoError(t, err)

		cookie := &http.Cookie{
			Name:  "user_id",
			Value: encoded,
		}

		isInvalid := IsInvalidUserID(cookie)
		assert.True(t, isInvalid)
	})

	t.Run("Tampered cookie (invalid signature)", func(t *testing.T) {
		userID := GenerateUserID()

		encoded, err := EncodeUserID(userID)
		require.NoError(t, err)

		// Подделываем cookie: меняем подпись
		tampered := encoded + "tampered"

		cookie := &http.Cookie{
			Name:  "user_id",
			Value: tampered,
		}

		isInvalid := IsInvalidUserID(cookie)
		assert.True(t, isInvalid, "tampered cookie must be marked invalid")
	})

	t.Run("Corrupted cookie value", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "user_id",
			Value: "%%%invalid%%%base64%%%==",
		}

		isInvalid := IsInvalidUserID(cookie)
		assert.True(t, isInvalid, "corrupted cookie must be marked invalid")
	})

	t.Run("Empty cookie value", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "user_id",
			Value: "",
		}

		isInvalid := IsInvalidUserID(cookie)
		assert.False(t, isInvalid, "empty cookie must be marked invalid")
	})
}
