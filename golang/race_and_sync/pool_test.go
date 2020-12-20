package race_and_sync

import "testing"

func BenchmarkWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithoutPool()
	}
}

func BenchmarkWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithPool()
	}
}
