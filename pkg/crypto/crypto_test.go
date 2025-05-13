package crypto

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
