package crypto

import "testing"

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
