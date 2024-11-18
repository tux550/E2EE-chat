package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	X3DHCore "tux.tech/x3dh/core"
)

type Server struct {
	docdb *DocDBManager
	memdb *MemDBManager
	apim  *AWSAPIManager
}

// Check if mock variable is set and if it is true
func IsMockSet() bool {
	mock := os.Getenv("MOCK")
	return mock == "true"
}

func NewServer() *Server {
	var docdb *DocDBManager
	var memdb *MemDBManager
	var apim *AWSAPIManager

	if IsMockSet() {
		// TEMPORARY NIL FOR LOCAL TESTING
		// Log mock mode
		log.Println("Running in mock mode")
		docdb = NewDocDBManager()
		memdb = nil
		apim = nil

	} else {
		// Log production mode
		log.Println("Running in production mode")
		docdb = NewDocDBManager()
		memdb = NewMemDBManager()
		apim = NewAWSAPIManager()

	}

	return &Server{
		docdb: docdb,
		memdb: memdb,
		apim:  apim,
	}
}

func (s *Server) RegisterClient(clientID string, bundle X3DHCore.X3DHClientBundle) error {
	// Register client in documentDB
	data := NewClientData(bundle)
	err := s.docdb.ClientUpsert(clientID, *data)
	return err
}

func (s *Server) IsClientRegistered(clientID string) (bool, error) {
	exists, err := s.docdb.ClientExists(clientID)
	return exists, err
}

func (s *Server) GetRemainingOTPCount(clientID string) (int, error) {
	return s.docdb.ClientOTPCount(clientID)
}

func (s *Server) ExpandOTPSet(clientID string, otps []X3DHCore.X3DHPublicOTP) error {
	return s.docdb.ClienOTPAppend(clientID, otps)
}

func (s *Server) GetClientBundle(clientID string) (X3DHCore.X3DHKeyBundle, bool, error) {
	clientData, err := s.docdb.ClienOTPPop(clientID, 1)
	if err != nil {
		return X3DHCore.X3DHKeyBundle{}, false, err
	}
	if len(clientData.Bundle.OtpSet) == 0 {
		return X3DHCore.X3DHKeyBundle{}, false, nil
	}
	otp := clientData.Bundle.OtpSet[0]

	return X3DHCore.X3DHKeyBundle{
		IK:  clientData.Bundle.IK,
		SPK: clientData.Bundle.SPK,
		OTP: otp,
	}, true, nil
}

func (s *Server) SendMessage(recipientID string, senderID string, msg X3DHCore.InitialMessage) bool {
	err := s.docdb.ClientMessageAppend(recipientID,
		MessageData{
			SenderID: senderID,
			Message:  msg,
		},
	)
	return err == nil
}

func (s *Server) GetMessage(clientID string) (MessageData, bool, error) {
	messageQueue, err := s.docdb.ClientMessagePop(clientID, 1)
	if err != nil {
		return MessageData{}, false, err
	}
	if len(messageQueue) == 0 {
		return MessageData{}, false, nil
	}
	msg := messageQueue[0]
	return msg, true, nil
}

func (s *Server) getMockConnectionEntry(connectionID string) (ConnectionEntry, error) {
	// If mock, request from mock-gateway:8082/connections/{connectionID}
	url := fmt.Sprintf("http://mock-gateway:8082/connections/%s", connectionID)
	resp, err := http.Get(url)

	entry := ConnectionEntry{}

	if err != nil {
		return entry, err
	}
	defer resp.Body.Close()
	// Parse response
	err = json.NewDecoder(resp.Body).Decode(&entry)
	if err != nil {
		return entry, err
	}
	return entry, nil
}

func (s *Server) getConnectionEntry(connectionID string) (ConnectionEntry, error) {

	if IsMockSet() {
		return s.getMockConnectionEntry(connectionID)
	} else {
		// Request from DynamoDB
		connection, err := s.memdb.GetConnection(connectionID)
		if err != nil {
			return ConnectionEntry{}, err
		}
		return *connection, nil
	}
}

func (s *Server) getMockConnectionEntriesByUser(userID string) ([]ConnectionEntry, error) {
	// If mock, request from mock-gateway:8082/username/{userID}
	url := fmt.Sprintf("http://mock-gateway:8082/username/%s", userID)
	resp, err := http.Get(url)

	entries := []ConnectionEntry{}

	if err != nil {
		return entries, err
	}
	defer resp.Body.Close()
	// Parse response
	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		return entries, err
	}
	return entries, nil
}

func (s *Server) getConnectionEntriesByUser(userID string) ([]ConnectionEntry, error) {
	if IsMockSet() {
		return s.getMockConnectionEntriesByUser(userID)
	} else {
		// Request from DynamoDB
		connections, err := s.memdb.GetUserConnections(userID)
		if err != nil {
			return []ConnectionEntry{}, err
		}
		// From list of pointers to list of values
		entries := make([]ConnectionEntry, len(connections))
		for i, conn := range connections {
			entries[i] = *conn
		}
		return entries, nil
	}
}

func (s *Server) mockWebSocketSendConnection(connectionID string, msg []byte) {
	// If mock, send to mock-gateway:8082/send/{connectionID}
	url := fmt.Sprintf("http://mock-gateway:8082/send/%s", connectionID)
	// Send message
	_, err := http.Post(url, "application/json", bytes.NewBuffer(msg))
	if err != nil {
		log.Printf("Failed to send message to connection %s: %v", connectionID, err)
	}
}

func (s *Server) WebSocketSendConnection(connectionID string, msg []byte) {
	// Send message to connection
	if IsMockSet() {
		s.mockWebSocketSendConnection(connectionID, msg)
	} else {
		// Send to connection
		log.Printf("WS send to connection %s: %s", connectionID, string(msg))
		s.apim.PostToConnection(connectionID, msg)
	}
}

func (s *Server) WebSocketSendUser(user string, msg []byte) {
	// Get connections by user
	var connections []ConnectionEntry
	connections, err := s.getConnectionEntriesByUser(user)
	if err != nil {
		log.Printf("Failed to get connections for user %s: %v", user, err)
		return
	}
	// For each connection, send message
	for _, conn := range connections {
		s.WebSocketSendConnection(conn.ConnectionID, msg)
	}

}
