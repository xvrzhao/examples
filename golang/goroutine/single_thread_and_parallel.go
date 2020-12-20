package goroutine

var data []int

const (
	countOfElement        = 1000000
	countOfParallel       = 25
	countOfSegmentElement = countOfElement / countOfParallel
)

// SingleThread calculates the sum in a single thread.
//
// Run the benchmark to compare with the Parallel function:
//   $ go test -bench=. -benchmem ./goroutine
//   BenchmarkSingleThread-8             1486            691141 ns/op               0 B/op          0 allocs/op
//   BenchmarkParallel-8                 4402            271651 ns/op             320 B/op          1 allocs/op
func SingleThread() int {
	var sum int
	for _, v := range data {
		sum += v
	}
	return sum
}

// Parallel calculates the sum in multiple goroutines (multiple threads).
func Parallel() int {
	ch := make(chan int, countOfParallel)
	for i := 0; i < countOfParallel; i++ {
		go func(i int) {
			var res int
			init := countOfSegmentElement * i
			for j := init; j < init+countOfSegmentElement; j++ {
				res += data[j]
			}
			ch <- res
		}(i)
	}
	var sum int
	for i := 0; i < countOfParallel; i++ {
		sum += <-ch
	}
	return sum
}

func init() {
	data = make([]int, countOfElement)
	for i := 0; i < countOfElement; i++ {
		data[i] = i
	}
}
