package x3dh_server

import (
	X3DHCore "tux.tech/x3dh/core"
)

type MessageData struct {
	SenderID string
	Message  X3DHCore.InitialMessage
}

type ClientData struct {
	// Bundle
	Bundle X3DHCore.X3DHClientBundle
	// Queue
	Queue []MessageData
}

type Server struct {
	clients map[string]*ClientData
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]*ClientData),
	}
}

func NewClientData(bundle X3DHCore.X3DHClientBundle) *ClientData {
	return &ClientData{
		Bundle: bundle,
		Queue:  make([]MessageData, 0),
	}
}

func (s *Server) RegisterClient(clientID string, bundle X3DHCore.X3DHClientBundle) {
	s.clients[clientID] = NewClientData(bundle)
}

func (s *Server) IsClientRegistered(clientID string) bool {
	_, ok := s.clients[clientID]
	return ok
}

func (s *Server) GetRemainingOTPCount(clientID string) int {
	c, ok := s.clients[clientID]
	if !ok {
		return 0
	}
	return len(c.Bundle.OtpSet)
}

func (s *Server) ExpandOTPSet(clientID string, otps []X3DHCore.X3DHPublicOTP) {
	c, ok := s.clients[clientID]
	if !ok {
		return
	}
	c.Bundle.OtpSet = append(c.Bundle.OtpSet, otps...)
}

func (s *Server) GetClientBundle(clientID string) (X3DHCore.X3DHKeyBundle, bool) {
	c, ok := s.clients[clientID]
	if !ok {
		return X3DHCore.X3DHKeyBundle{}, false
	}
	// Build the key bundle
	// Identity Key
	ik := c.Bundle.IK
	// Signed Pre Key
	spk := c.Bundle.SPK
	// One Time Pre Key
	if len(c.Bundle.OtpSet) == 0 {
		return X3DHCore.X3DHKeyBundle{}, false
	}
	otp := c.Bundle.OtpSet[0]
	c.Bundle.OtpSet = c.Bundle.OtpSet[1:]
	// Return
	return X3DHCore.X3DHKeyBundle{
		IK:  ik,
		SPK: spk,
		OTP: otp,
	}, true
}

func (s *Server) SendMessage(recipientID string, senderID string, msg X3DHCore.InitialMessage) bool {
	c, ok := s.clients[recipientID]
	if !ok {
		return false
	}
	c.Queue = append(c.Queue, MessageData{
		SenderID: senderID,
		Message:  msg,
	})
	return true
}

func (s *Server) GetMessage(clientID string) (MessageData, bool) {
	c, ok := s.clients[clientID]
	if !ok {
		return MessageData{}, false
	}
	if len(c.Queue) == 0 {
		return MessageData{}, false
	}
	msg := c.Queue[0]
	c.Queue = c.Queue[1:]
	return msg, true
}
