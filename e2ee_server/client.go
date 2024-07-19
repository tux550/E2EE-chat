package main

import (
	"encoding/json"
	"fmt"

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
		if err != nil || mt == websocket.CloseMessage {
			break // Exit loop
		}
		if mt == websocket.TextMessage {
			client.HandleMessage(message)
		}
	}
	client.Disconnect()
}

func (client *WsClient) HandleMessage(rawMessage []byte) {
	// Log
	fmt.Println("Received message from", client.username, ":", string(rawMessage))
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
	case "status":
		client.HandleUserStatus(message.Params)
	case "upload_new_otps":
		client.HandleUploadNewOTPs(message.Params)
	}
}

func buildOutboundMessage(params interface{}, method string) (*api.OutboundMessage, error) {
	// Marshal
	marshalledParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	// Build API call
	api_call := &api.OutboundMessage{
		Method: method,
		Params: marshalledParams,
	}
	return api_call, nil
}

func (client *WsClient) Disconnect() {
	client.server.UnsetClient(client)
	client.conn.Close()
	close(client.send) // Close write pump
}

func getLowOTPNotification() []byte {
	notification, err := buildOutboundMessage(&api.NotifyLowOTP{}, "notify_low_otp")
	if err != nil {
		fmt.Println("Error marshalling notification to notify_low_otp")
		return nil
	}
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshalling notification to notify_low_otp")
		return nil
	}
	return notificationBytes
}

func (client *WsClient) HandleUploadNewOTPs(rawParams json.RawMessage) {
	params := &api.RequestUploadOTPs{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// Register OTPs
	client.server.X3DHServer.ExpandOTPSet(client.username, params.OTPs)
	fmt.Println("User", client.username, "uploaded #", len(params.OTPs), "new OTPs")
}

func (client *WsClient) HandleGetUserBundle(rawParams json.RawMessage) {
	params := &api.RequestUserBundle{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}

	// Get the bundle
	bundle, ok, err := client.server.X3DHServer.GetClientBundle(params.UserID)
	if err != nil {
		fmt.Println("Db error getting bundle for user", params.UserID)
		return
	}
	fmt.Println("User", client.username, "requested bundle for user", params.UserID, ":", ok)

	// Notify recipient if otp is running low
	count, err := client.server.X3DHServer.GetRemainingOTPCount(params.UserID)
	if err != nil {
		fmt.Println("Error getting remaining OTP count for user", params.UserID)
		return
	}
	if count < 3 {
		// Notify user
		fmt.Println("Notifying user", params.UserID, "that OTP is running low")
		// Send notification
		notificationBytes := getLowOTPNotification()
		client.server.SendNotificationToUser(params.UserID, notificationBytes)
	}
	// Send response
	response, err := buildOutboundMessage(&api.ResponseUserBundle{
		Success: ok,
		Bundle:  bundle,
	}, "get_bundle")
	if err != nil {
		fmt.Println("Error marshalling response to get_bundle")
		return
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling fail response to get_bundle")
		return
	}
	client.send <- responseBytes

}

func (client *WsClient) HandleUploadBundle(rawParams json.RawMessage) {
	params := &api.RequestUploadBundle{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// Check if the user is the same
	if client.username != params.UserID {
		fmt.Println("User", client.username, "attempted to upload bundle for user", params.UserID)

		response, err := buildOutboundMessage(&api.ResponseUploadBundle{
			Success: false,
		}, "upload_bundle")
		if err != nil {
			fmt.Println("Error marshalling response to upload_bundle")
			return
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling fail response to upload_bundle")
			return
		}
		client.send <- responseBytes
	}
	// Register
	client.server.X3DHServer.RegisterClient(params.UserID, params.Bundle)
	fmt.Println("User", client.username, "uploaded bundle for user", params.UserID)
	// Send response
	response, err := buildOutboundMessage(&api.ResponseUploadBundle{
		Success: true,
	}, "upload_bundle")
	if err != nil {
		fmt.Println("Error marshalling response to upload_bundle")
		return
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling success response to upload_bundle")
		return
	}
	client.send <- responseBytes
}

func (client *WsClient) HandleSendMessage(rawParams json.RawMessage) {
	params := &api.RequestSendMsg{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// Send message
	ok := client.server.X3DHServer.SendMessage(params.RecipientID, client.username, params.MessageData)
	fmt.Println("User", client.username, "sent message to user", params.RecipientID, ":", ok)
	// Send response
	response, err := buildOutboundMessage(&api.ResponseSendMsg{
		Success: ok,
	}, "send_message")
	if err != nil {
		fmt.Println("Error marshalling response to send_message")
		return
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling success response to send_message")
		return
	}
	client.send <- responseBytes
	// Notify recipient of message only if successful
	if !ok {
		return
	}
	// Send notification
	notification, err := buildOutboundMessage(&api.NotifyNewMessage{
		SenderID: client.username,
	}, "notify_new_message")
	if err != nil {
		fmt.Println("Error marshalling notification to notify_new_message")
		return
	}
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshalling notification to notify_new_message")
		return
	}
	client.server.SendNotificationToUser(params.RecipientID, notificationBytes)
}

func (client *WsClient) HandleReceiveMessage(rawParams json.RawMessage) {
	params := &api.RequestReceiveMsg{}
	err := json.Unmarshal(rawParams, params)
	if err != nil {
		return
	}
	// Unqueue message
	messageData, ok, err := client.server.X3DHServer.GetMessage(client.username)
	if err != nil {
		fmt.Println("Error getting message for user", client.username)
		return
	}
	if !ok {
		fmt.Println("User", client.username, "requested message but none available")
		// Send response
		response, err := buildOutboundMessage(&api.ResponseReceiveMsg{
			Success: false,
		}, "receive_message")
		if err != nil {
			fmt.Println("Error marshalling response to receive_message")
			return
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling fail response to receive_message")
			return
		}
		client.send <- responseBytes
		return
	}
	fmt.Println("User", client.username, "received message from user", messageData.SenderID)
	// Send response
	response, err := buildOutboundMessage(&api.ResponseReceiveMsg{
		Success:     true,
		SenderID:    messageData.SenderID,
		MessageData: messageData.Message,
	}, "receive_message")
	if err != nil {
		fmt.Println("Error marshalling response to receive_message")
		return
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling success response to receive_message")
		return
	}
	client.send <- responseBytes
}

func (client *WsClient) HandleUserStatus(rawParams json.RawMessage) {
	// Check if the user is registered
	//registered := client.server.X3DHServer.IsClientRegistered(params.UserID)
	//fmt.Println("User", client.username, "checked if user", params.UserID, "is registered")

	// Check if self is registered
	registered, err := client.server.X3DHServer.IsClientRegistered(client.username)
	if err != nil {
		fmt.Println("Error checking if user", client.username, "is registered")
		return
	}
	fmt.Println("User", client.username, "checked if self is registered")

	// Send response
	response, err := buildOutboundMessage(&api.ResponseUserStatus{
		Success: registered,
	}, "status")
	if err != nil {
		fmt.Println("Error marshalling response to status")
		return
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling success response to status")
		return
	}
	client.send <- responseBytes
	// If not registered, do not notify otp
	if !registered {
		return
	}
	// Notify recipient if otp is running low
	count, err := client.server.X3DHServer.GetRemainingOTPCount(client.username)
	if err != nil {
		fmt.Println("Error getting remaining OTP count for user", client.username)
		return
	}
	if count < 3 {
		// Notify user
		fmt.Println("Notifying user", client.username, "that OTP is running low")
		// Send notification
		notificationBytes := getLowOTPNotification()
		client.send <- notificationBytes
	}
}
