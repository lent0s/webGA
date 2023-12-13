package cmd

import (
	"context"
	"github.com/lent0s/webGA/logger"
	"log"
	"net/http"
	"sync"
)

func ServerUp(exit chan struct{}, w *sync.WaitGroup) {

	defer w.Done()

	p := getDefaultPage()

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.home)
	mux.HandleFunc("/getSettings", getSettings)
	mux.HandleFunc("/getStateInstance", getStateInstance)
	mux.HandleFunc("/sendMessage", sendMessage)
	mux.HandleFunc("/sendFileByUrl", sendFileByUrl)

	logs := logger.InitLogger()
	log.SetOutput(logs.Log2Way)
	addr := getHost()
	server := &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: log.New(logs.Log2Way, "SERVER: ", 0),
	}
	logs.Log2Way.Info().Msgf("server is running on %s", addr)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logs.Log2Way.Fatal().Msgf("%s", err)
		}
		logs.Log2Way.Info().Msgf("server close on %s", addr)
	}()

	<-exit
	if err := server.Shutdown(context.Background()); err != nil {
		logs.Log2Way.Error().Msgf("on shutdown server: %s", err)
	}
	wg.Wait()
}
