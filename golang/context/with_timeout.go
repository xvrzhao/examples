package context

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RunTimeoutExample 演示了如何通过 context、select 和 waitgroup 机制
// 实现子协程派生与关闭。
//
//	context: 实现了从主协程到子协程的关闭信息的传递
//	select: 实现了子协程内部对父协程关闭信号的监控
//	waitgroup: 实现了主协程对子协程是否退出的感知
func RunTimeoutExample() {
	timeout := time.Second * 3
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

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
					fmt.Printf("returned: goroutine 2, live %v\n", time.Since(start))
					return
				default:
				}

				// working
				time.Sleep(time.Second / 2)
			}
		}(ctx, wg, start)

		var second int
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("returned: goroutine 1, live %v\n", time.Since(start))
				return
			default:
			}

			// working
			s := time.Now().Second()
			if second != s {
				second = s
				fmt.Printf("now, second: %ds\n", s)
			}
		}
	}(ctx, &wg, start)

	wg.Wait()
}
