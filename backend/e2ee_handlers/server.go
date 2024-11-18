package main

import (
	"log"

	X3DHCore "tux.tech/x3dh/core"
)

type Server struct {
	docdb *DocDBManager
	memdb *MemDBManager
	apim  *AWSAPIManager
}

func NewServer() *Server {

	// TEMPORARY NIL FOR LOCAL TESTING
	// Create db managers
	docdb := NewDocDBManager()
	var memdb *MemDBManager = nil //NewMemDBManager()
	var apim *AWSAPIManager = nil //NewAWSAPIManager()

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

func (s *Server) getConnectionEntry(connectionID string) (ConnectionEntry, error) {
	// Mock request connection from DynamoDB
	return ConnectionEntry{
		ConnectionID: connectionID,
		Username:     "user1",
	}, nil
}

func (s *Server) WebSocketSendConnection(connectionID string, msg []byte) {
	// Mock
	log.Printf("WS send to connection %s: %s", connectionID, string(msg))
}

func (s *Server) WebSocketSendUser(user string, msg []byte) {
	// Mock
	// Should find the connectionID of the user and send the message
	log.Printf("WS send to user %s: %s", user, string(msg))
}
