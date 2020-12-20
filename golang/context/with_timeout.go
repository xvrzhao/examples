package context

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func RunTimeoutExample() {
	timeout := time.Second * 3
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context, wg *sync.WaitGroup, start time.Time) {
		defer wg.Done()
		fmt.Println("started: goroutine 1")

		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, start time.Time) {
			defer wg.Done()
			fmt.Println("started: goroutine 2 derived by 1")
			for {
				select {
				case <-ctx.Done():
					fmt.Println("returned: goroutine 2, live", time.Since(start))
					return
				default:
				}
			}
		}(ctx, wg, start)

		var second int
		for {
			select {
			case <-ctx.Done():
				fmt.Println("returned: goroutine 1, live", time.Since(start))
				return
			default:
				s := time.Now().Second()
				if second != s {
					second = s
					fmt.Printf("now, second: %ds\n", s)
				}
			}
		}
	}(ctx, &wg, start)

	wg.Wait()
}
