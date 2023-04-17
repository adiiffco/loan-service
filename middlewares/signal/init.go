package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var cleanup = make(chan int)

func multiSignalHandler(signal os.Signal) {
	switch signal {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
		fmt.Println("Signal:", signal.String())
		cleanup <- 1
	default:
		fmt.Println("Unhandled/unknown signal")
	}
}

func SetupSignals() {
	fmt.Println("Initializing signals")
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)

	go func() {
		for {
			s := <-sigchnl
			multiSignalHandler(s)
		}
	}()
}

func CleanupOnSignal(f func()) {
	go func() {
		<-cleanup
		close(cleanup)
		f()
	}()
}
