package main

import (
	"log"
	config "server/config"
)

func main() {
	// Start the server
	err := config.StartServer()
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
