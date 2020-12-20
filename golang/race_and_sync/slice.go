package race_and_sync

import (
	"fmt"
	"sync"
)

// SliceSafe demonstrates the concurrency-safe scenario of slice.
// Run the following command under the project root directory to test race:
//   $ go test -race -run=SliceSafe -v ./race_and_sync
func SliceSafe() {
	s := make([]int, 10)
	wg := new(sync.WaitGroup)
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			s[i] = i * i
		}(i)
	}

	wg.Wait()
	fmt.Println(s)
}

// SliceUnsafe demonstrates the non-concurrency-safe scenario of slice.
// Run the following command under the project root directory to test race:
//   $ go test -run=SliceUnsafe -v -race -count=1 ./race_and_sync
func SliceUnsafe() {
	s := make([]int, 0)
	wg := new(sync.WaitGroup)
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			s = append(s, i)
		}(i)
	}

	wg.Wait()
	fmt.Println(s)
}
