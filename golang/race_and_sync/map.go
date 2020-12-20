package race_and_sync

import (
	"fmt"
	"sync"
)

// MapIsNonConcurrentSafe demonstrates that a plain map is not concurrency-safe.
// Run the following command under the project root directory to test race:
//   $ go test -race -run=MapIsNonConcurrencySafe ./race_and_sync
func MapIsNonConcurrencySafe() {
	m := make(map[int]int)
	done := make(chan bool)

	go func() {
		m[1] = 1 // First conflicting access.
		done <- true
	}()

	m[0] = 0 // // Second conflicting access.
	<-done
}

type myMap struct {
	sync.RWMutex
	m map[int]int
}

// ConcurrencySafeMapWithLock demonstrates that a map with the lock is concurrency-safe.
// Run the following command under the project root directory to test race:
//   $ go test -race -run=ConcurrencySafeMapWithLock ./race_and_sync
func ConcurrencySafeMapWithLock() {
	mm := new(myMap)
	mm.m = make(map[int]int)

	wg := new(sync.WaitGroup)
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()

			mm.Lock()
			mm.m[i] = i * 2
			mm.Unlock()

			mm.RLock()
			fmt.Printf("%d => %d\n", i, mm.m[i])
			mm.RUnlock()
		}(i)
	}

	wg.Wait()

	fmt.Println("Map:")
	for k, v := range mm.m {
		fmt.Printf("\t%d => %d\n", k, v)
	}
}

// SyncMapIsConcurrencySafe demonstrates that a sync.Map is concurrency-safe.
// Run the following command under the project root directory to test race:
//   $ go test -race -run=SyncMapIsConcurrencySafe ./race_and_sync
func SyncMapIsConcurrencySafe() {
	m := new(sync.Map)
	done := make(chan bool)

	go func() {
		m.Store("red", "apple")
		done <- true
	}()

	m.Store("yellow", "banana")
	<-done

	m.Range(func(k, v interface{}) bool {
		fmt.Println(k.(string), v.(string))
		return true
	})
}
