package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	api "tux.tech/e2ee/api"
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
		client.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (client *WsClient) ReadPump() {
	for {
		mt, message, err := client.conn.ReadMessage()
		if err != nil || mt != websocket.CloseMessage {
			break // Exit loop
		}
		if mt == websocket.TextMessage {
			client.HandleMessage(message)
		}
	}
	client.Disconnect()
}

func (client *WsClient) HandleMessage(rawMessage []byte) {
	// Parse JSON
	message := &api.InboundMessage{}
	err := json.Unmarshal(rawMessage, message)
	if err != nil {
		return
	}
	switch message.Method {
	case "get_bundle":
		client.HandleGetUserBundle(message.Params)
	case "upload_bundle":
		client.HandleUploadBundle(message.Params)
	case "send_message":
		client.HandleSendMessage(message.Params)
	case "receive_message":
		client.HandleReceiveMessage(message.Params)
	}

}

func (client *WsClient) Disconnect() {
	client.server.UnsetClient(client)
	client.conn.Close()
	close(client.send) // Close write pump
}

func (client *WsClient) HandleGetUserBundle(rawParams json.RawMessage) {
	params := &api.RequestUserBundle{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// TODO: Implement
}

func (client *WsClient) HandleUploadBundle(rawParams json.RawMessage) {
	params := &api.RequestUploadBundle{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// TODO: Implement
}

func (client *WsClient) HandleSendMessage(rawParams json.RawMessage) {
	params := &api.RequestSendMsg{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// TODO: Implement
}

func (client *WsClient) HandleReceiveMessage(rawParams json.RawMessage) {
	params := &api.RequestReceiveMsg{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// TODO: Implement
}
