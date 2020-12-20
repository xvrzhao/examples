package channel

import (
	"fmt"
	"sync"
)

// CloseAChannel demonstrates that the close function is always only used by a producer,
// and when a channel is closed, the for range loop of that channel will be stopped.
func CloseAChannel() {
	ch := make(chan string)

	var wg sync.WaitGroup
	wg.Add(3)

	// producer
	go func(ch chan<- string, wg *sync.WaitGroup) {
		defer wg.Done()
		ch <- "apple"
		ch <- "orange"
		ch <- "banana"
		close(ch)
	}(ch, &wg)

	// consumer
	go func(ch <-chan string, wg *sync.WaitGroup) {
		defer wg.Done()
		// when ch is closed, the for loop will stop
		for item := range ch {
			fmt.Printf("consumer 1 received: %s\n", item)
		}
		fmt.Println("consumer 1 received the close signal from producer.")
	}(ch, &wg)

	go func(ch <-chan string, wg *sync.WaitGroup) {
		defer wg.Done()
		for item := range ch {
			fmt.Printf("consumer 2 received: %s\n", item)
		}
		fmt.Println("consumer 2 received the close signal from producer.")
	}(ch, &wg)

	wg.Wait()
}

// RetrieveItemFromClosedChannel demonstrates that when receiving a value
// from a closed channel, a zero-value is received.
func RetrieveItemFromClosedChannel() {
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan int)

	go func(ch chan<- int) {
		defer wg.Done()
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch) // close the channel when number 1~3 have been sent to it
	}(ch)

	go func(ch <-chan int) {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			n, ok := <-ch
			fmt.Printf("%d %v\n", n, ok)
			// 1 true
			// 2 true
			// 3 true
			// 0 false
			// 0 false
		}
	}(ch)

	wg.Wait()
}
