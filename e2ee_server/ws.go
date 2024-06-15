package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsServer struct {
	clients map[*WsClient]bool
	mu      sync.Mutex
}

func NewWsServer() *WsServer {
	return &WsServer{
		clients: make(map[*WsClient]bool),
	}
}

func (server *WsServer) SetClient(client *WsClient) {
	server.mu.Lock()
	server.clients[client] = true
	server.mu.Unlock()
}

func (server *WsServer) UnsetClient(client *WsClient) {
	server.mu.Lock()
	delete(server.clients, client)
	server.mu.Unlock()
}

func (server *WsServer) connnect(w http.ResponseWriter, r *http.Request) {
	// Get user fom header
	user := r.Header.Get("User")
	if user == "" {
		http.Error(w, "No user", http.StatusBadRequest)
		return
	}
	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	client := NewWsClient(user, server, conn)

	server.SetClient(client)

	fmt.Println("New connection from", user)
	go client.WritePump()
	go client.ReadPump()
}
