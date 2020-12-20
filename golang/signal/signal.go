package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func RunReceiveInterruptSignal() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT) // if not specify SIGINT, all signals will be received into signalChannel

	done := make(chan int)
	go func(done chan<- int) {
		sig := <-signalChannel
		fmt.Printf("receive signal: %v\n", sig)
		done <- 1
	}(done)

	fmt.Println("waiting signal...")
	<-done
}
