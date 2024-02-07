package cmd

import (
	"context"
	"github.com/lent0s/webGA/link"
	"github.com/lent0s/webGA/logger"
	"log"
	"net/http"
	"sync"
)

func ServerUp(exit chan struct{}, wg *sync.WaitGroup, logs *logger.Logs) {

	defer wg.Done()

	p := getDefaultPage()

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.home)

	//mux.HandleFunc("/", p.web)

	//mux.HandleFunc("/", p.warm)

	mux.HandleFunc("/link", link.Link)

	mux.HandleFunc("/getSettings", getSettings)
	mux.HandleFunc("/getStateInstance", getStateInstance)
	mux.HandleFunc("/sendMessage", sendMessage)
	mux.HandleFunc("/sendFileByUrl", sendFileByUrl)

	addr := getHost(logs)
	server := &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: log.New(logs.Log2Way, "SERVER: ", 0),
	}
	logs.Log2Way.Info().Msgf("SERVER: running on %s", addr)

	wgServer := sync.WaitGroup{}
	wgServer.Add(1)
	go func() {
		defer wgServer.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logs.Log2Way.Fatal().Msgf("%s", err)
		}
		logs.Log2Way.Info().Msgf("SERVER: close on %s", addr)
	}()

	<-exit
	if err := server.Shutdown(context.Background()); err != nil {
		logs.Log2Way.Error().Msgf("on shutdown server: %s", err)
	}
	wgServer.Wait()
}
