package main

import (
	"github.com/lent0s/webGA/cmd"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	exit := make(chan struct{})
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal,
		syscall.SIGTERM,
		syscall.SIGINT)
	go waitExit(exitSignal, exit)

	wg := sync.WaitGroup{}
	wg.Add(1)
	cmd.ServerUp(exit, &wg)

	wg.Wait()
}

func waitExit(exitSignal chan os.Signal, exit chan struct{}) {

	<-exitSignal
	close(exit)
}
