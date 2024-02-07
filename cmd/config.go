package cmd

import (
	"github.com/lent0s/webGA/logger"
	"os"
)

func getHost(logs *logger.Logs) string {

	port := os.Getenv("PORT")
	if port == "" {
		port = "9012"
		logs.Log2Way.Info().Msgf("$PORT must be set - using default: %s", port)
	}

	return ":" + port
}
