package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	api "tux.tech/e2ee/api"
	x3dh_core "tux.tech/x3dh/core"
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
	// Parse JSON
	params := &api.RequestUserBundle{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// Get user bundle
	bundle, ok := client.server.x3dh_server.GetClientBundle(params.UserID)
	if !ok {
		return
	}
	// To API response
	response := &api.OutboundMessage{
		Method: "get_bundle",
		Params: api.ResponseUserBundle(bundle),
	}
	// Send response
	responseJson, err := json.Marshal(response)
	if err != nil {
		return
	}
	client.send <- responseJson
}

func (client *WsClient) HandleUploadBundle(rawParams json.RawMessage) {
	// Parse JSON
	params := &api.RequestUploadBundle{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// Register client
	client.server.x3dh_server.RegisterClient(client.username, *(*x3dh_core.X3DHClientBundle)(params))
	// To API response
	response := &api.OutboundMessage{
		Method: "upload_bundle",
		Params: api.ResponseUploadBundle{Success: true},
	}
	// Send response
	responseJson, err := json.Marshal(response)
	if err != nil {
		return
	}
	client.send <- responseJson
}

func (client *WsClient) HandleSendMessage(rawParams json.RawMessage) {
	// Parse JSON
	params := &api.RequestSendMsg{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// Send message
	client.server.x3dh_server.SendMessage(client.username, params.RecipientID, params.MessageData)
	// To API response
	response := &api.OutboundMessage{
		Method: "send_message",
		Params: api.ResponseSendMsg{Success: true},
	}
	// Send response
	responseJson, err := json.Marshal(response)
	if err != nil {
		return
	}
	client.send <- responseJson
}

func (client *WsClient) HandleReceiveMessage(rawParams json.RawMessage) {
	// Parse JSON
	params := &api.RequestReceiveMsg{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// Receive message
	msg, ok := client.server.x3dh_server.GetMessage(client.username)
	if !ok {
		return
	}
	// To API response
	response := &api.OutboundMessage{
		Method: "receive_message",
		Params: api.ResponseReceiveMsg{
			SenderID:    msg.SenderID,
			RecipientID: msg.RecipientID,
			MessageData: msg.MessageData,
		},
	}
	// Send response
	responseJson, err := json.Marshal(response)
	if err != nil {
		return
	}
	client.send <- responseJson
}
