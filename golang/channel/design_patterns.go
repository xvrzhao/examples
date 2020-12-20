package channel

import (
	"fmt"
	"sync"
)

// FactoryPattern demonstrates the factory pattern of channel.
// The factory function returns a channel that produces data.
func FactoryPattern() {
	ch := factory()
	for n := range ch {
		fmt.Println(n)
	}
}

func factory() <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

// FindPrimeNumbersBySieves finds prime numbers by multiple sieves, see picture:
// https://raw.githubusercontent.com/unknwon/the-way-to-go_ZH_CN/master/eBook/images/14.2_fig14.2.png
func FindPrimeNumbersBySieves() {
	ch := make(chan int)
	go func(out chan<- int) {
		for i := 2; i <= 1000; i++ {
			out <- i
		}
		close(out)
	}(ch)

	var wg sync.WaitGroup

	for {
		if prime, ok := <-ch; ok {
			fmt.Printf("%d ", prime)
			ch1 := make(chan int)
			wg.Add(1)
			go func(in, out chan int, prime int) {
				defer wg.Done()
				filter(in, out, prime)
			}(ch, ch1, prime)
			ch = ch1
		} else {
			break
		}
	}

	wg.Wait()
}

func filter(in, out chan int, prime int) {
	for i := range in {
		if i%prime != 0 {
			out <- i
		}
	}
	close(out)
}

// FindPrimeNumbersBySievesFac is factory-pattern version of FindPrimeNumbersBySieves.
func FindPrimeNumbersBySievesFac() {
	for prime := range sieveFac() {
		fmt.Printf("%d ", prime)
	}
}

func sieveFac() (primes chan int) {
	primes = make(chan int)
	go func() {
		defer close(primes)
		ch := generatorFac()
		for {
			if prime, ok := <-ch; ok {
				primes <- prime
				ch = filterFac(ch, prime)
			} else {
				break
			}
		}
	}()
	return
}

func generatorFac() (out chan int) {
	out = make(chan int)
	go func() {
		defer close(out)
		for i := 2; i < 1000; i++ {
			out <- i
		}
	}()
	return
}

func filterFac(in chan int, prime int) (out chan int) {
	out = make(chan int)
	go func() {
		defer close(out)
		for i := range in {
			if i%prime != 0 {
				out <- i
			}
		}
	}()
	return
}
