package channel

import (
	"sync"
	"testing"
)

func calculateWithoutMutex() (sum int) {
	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			sum += 1
			wg.Done()
		}()
	}

	wg.Wait()
	return
}

func calculateWithMutex() (sum int) {
	m := NewMutex()

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			m.Lock()
			sum += 1
			m.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	return
}

func TestMutex(t *testing.T) {
	for i := 0; i < 5; i++ {
		sum1 := calculateWithoutMutex()
		sum2 := calculateWithMutex()
		t.Logf("without mutex: %d, with mutex: %d\n", sum1, sum2)
	}
}
