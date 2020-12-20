package channel

type mutex struct {
	ch chan interface{}
}

// NewMutex return a mutex that was implemented through a buffer channel.
//
// See the usage of the mutex in ./mutex_test.go, and run the unit test which compared sum calculations with and without mutex.
//   $ go test -v -run=Mutex -count=1  ./channel
func NewMutex() *mutex {
	return &mutex{ch: make(chan interface{}, 1)}
}

func (m *mutex) Lock() {
	m.ch <- new(interface{})
}

func (m *mutex) Unlock() {
	if m.ch == nil || len(m.ch) == 0 {
		panic("unlock on an unlocked mutex")
	}
	<-m.ch
}
