package main

import (
	"log"
	"net/http"
)

// Websocket - ECS
var (
	appVersion = "1.0.2"
)

func init() {
	// Log version
	log.Println("Version:", appVersion)
}

func main() {
	// Initialize API server
	server := NewServer()
	// Initialize API handler
	router := NewAPIHandler(server)
	// Start API server
	err := http.ListenAndServe(":8080", http.HandlerFunc(router.HandleRequests))
	if err != nil {
		log.Fatal(err)
	}
}
