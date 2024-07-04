package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	x3dh_server "tux.tech/x3dh/server"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsServer struct {
	clients     map[*WsClient]bool
	mu          sync.Mutex
	x3dh_server *x3dh_server.Server
}

func NewWsServer() *WsServer {
	return &WsServer{
		clients:     make(map[*WsClient]bool),
		x3dh_server: x3dh_server.NewServer(),
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
