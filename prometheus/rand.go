package main

import (
	"math/rand"
	"sync"
	"time"
)

var seedOnce sync.Once

// Seed do rand.Seed() Once.
func Seed() {
	seedOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})
}
