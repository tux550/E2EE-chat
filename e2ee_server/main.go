package main

import (
	"net/http"
)

func main() {
	server := NewWsServer()
	http.HandleFunc("/ws", server.connnect)
	http.ListenAndServe(":8765", nil)
}
