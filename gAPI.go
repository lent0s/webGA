package main

import (
	"github.com/lent0s/webGA/cmd"
	"github.com/lent0s/webGA/logger"
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

	logs := logger.InitLogger()
	cmd.ServerUp(exit, &wg, logs)

	wg.Wait()
	logger.StopLogger(logs)
}

func waitExit(exitSignal chan os.Signal, exit chan struct{}) {

	<-exitSignal
	close(exit)
}
