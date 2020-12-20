package race_and_sync

import (
	"fmt"
	"sync"
)

// NoCopy illustrates that some structs in go std lib annotates `must not be copied after first use`,
// since there cannot be two instances which contain certain pointer-type fields refer to the same
// memory at the same time, otherwise concurrency is not safe.
//
// So, instances of these structs cannot pass to function params by value, and cannot assign value to
// other variables, those are all copy options.
//
// Most of those structs contains the noCopy field, like sync.Mutex, sync.Pool, etc.
// So, `go vet` can check that whether struct instances are copied.
func NoCopy() {
	// 1. wrong demo, run `$ go vet ./race_and_sync` to check
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg sync.WaitGroup) {
		wg.Done()
	}(wg)
	wg.Wait()

	// 2. right demo
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func(wg *sync.WaitGroup) {
		wg.Done()
	}(&wg1)
	wg1.Wait()

	// 3. wrong demo, run `$ go vet ./race_and_sync` to check
	var wg2 sync.WaitGroup
	wg3 := wg2
	fmt.Println(wg3)

	// 4. self-made nocopy struct
	type Person struct {
		myNoCopy
		Name string
	}
	p := Person{Name: "xvr"}
	fmt.Println(p) // wrong, copied to param
}

type myNoCopy struct{}

func (*myNoCopy) Lock()   {}
func (*myNoCopy) Unlock() {}
