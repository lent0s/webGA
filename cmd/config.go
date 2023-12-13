package cmd

import (
	"log"
	"os"
)

func getHost() string {

	port := os.Getenv("PORT")
	if port == "" {
		port = "9012"
		log.Printf("$PORT must be set- Using default: %s\n", port)
	}

	return ":" + port
}
