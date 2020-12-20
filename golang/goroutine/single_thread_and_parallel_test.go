package goroutine

import "testing"

func BenchmarkSingleThread(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SingleThread()
	}
}

func BenchmarkParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parallel()
	}
}
