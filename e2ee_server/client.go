package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	username string
	server   *WsServer
	conn     *websocket.Conn
	send     chan []byte
}

func NewWsClient(username string, server *WsServer, conn *websocket.Conn) *WsClient {
	return &WsClient{
		username: username,
		server:   server,
		conn:     conn,
		send:     make(chan []byte),
	}
}

func (client *WsClient) WritePump() {
	for message := range client.send {
		client.conn.WriteMessage(websocket.BinaryMessage, message)
	}
}

func (client *WsClient) ReadPump() {
	for {
		mt, message, err := client.conn.ReadMessage()
		if err != nil || mt != websocket.CloseMessage {
			break // Exit loop
		}
		if mt == websocket.BinaryMessage {
			client.HandleMessage(message)
		}
	}
	client.Disconnect()
}

func (client *WsClient) HandleMessage(message []byte) {
	// TODO: Handle message
	fmt.Println("Received message:", message)
}

func (client *WsClient) Disconnect() {
	client.server.UnsetClient(client)
	client.conn.Close()
	close(client.send) // Close write pump
}
