package race_and_sync

import (
	"sync"
)

var (
	pool *sync.Pool
)

const (
	numOfGos   = 1 << 10 // 1024
	numOfBytes = 1 << 20 // 1M bytes
)

// Compare two following functions by command:
//   $ go test -bench=. -run=Pool -benchmem ./race_and_sync

func WithoutPool() {
	wg := new(sync.WaitGroup)
	wg.Add(numOfGos)
	for i := 0; i < numOfGos; i++ {
		go func() {
			bs := make([]byte, numOfBytes)
			for j := 0; j < len(bs); j++ {
				bs[j] = 0x00
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func WithPool() {
	wg := new(sync.WaitGroup)
	wg.Add(numOfGos)
	for i := 0; i < numOfGos; i++ {
		go func() {
			bs := pool.Get().(*[]byte)
			for j := 0; j < len(*bs); j++ {
				(*bs)[j] = 0x00
			}
			pool.Put(bs)
			wg.Done()
		}()
	}
	wg.Wait()
}

func init() {
	pool = &sync.Pool{New: func() interface{} {
		bs := make([]byte, numOfBytes)
		return &bs
	}}
}
