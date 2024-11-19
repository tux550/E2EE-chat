package main

import X3DHCore "tux.tech/x3dh/core"

// === DynamoDB ===
// Connection Tables
type ConnectionEntry struct {
	ConnectionID string `json:"connectionId"`
	Username     string `json:"username"`
}

// === DocumentDB ===
// Message Tables
type MessageData struct {
	SenderID string
	Message  X3DHCore.InitialMessage
}

// Client Tables
type ClientData struct {
	// Bundle
	Bundle X3DHCore.X3DHClientBundle
	// Queue
	Queue []MessageData
}

func NewClientData(bundle X3DHCore.X3DHClientBundle) *ClientData {
	return &ClientData{
		Bundle: bundle,
		Queue:  make([]MessageData, 0),
	}
}
