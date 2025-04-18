package crypto

import "testing"

func BenchmarkDecodeUserID(b *testing.B) {

	eUserId, _ := EncodeUserID("test")
	eIUserId, _ := EncodeUserIDRSA("test")

	b.Run("simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DecodeUserID(eUserId)
		}
	})

	b.Run("shuffle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DecodeUserIDRSA(eIUserId)
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
