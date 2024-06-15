package x3dh_server

import (
	"x3dh_core"
)

type ClientData struct {
	// Bundle
	Bundle x3dh_core.X3DHClientBundle
	// Queue
	Queue []x3dh_core.InitialMessage
}

type Server struct {
	clients map[string]*ClientData
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]*ClientData),
	}
}

func NewClientData(bundle x3dh_core.X3DHClientBundle) *ClientData {
	return &ClientData{
		Bundle: bundle,
		Queue:  make([]x3dh_core.InitialMessage, 0),
	}
}

func (s *Server) RegisterClient(clientID string, bundle x3dh_core.X3DHClientBundle) {
	s.clients[clientID] = NewClientData(bundle)
}

func (s *Server) GetClientBundle(clientID string) (x3dh_core.X3DHKeyBundle, bool) {
	c, ok := s.clients[clientID]
	if !ok {
		return x3dh_core.X3DHKeyBundle{}, false
	}
	// Build the key bundle
	// Identity Key
	ik := c.Bundle.IK
	// Signed Pre Key
	spk := c.Bundle.SPK
	// One Time Pre Key
	if len(c.Bundle.OtpSet) == 0 {
		return x3dh_core.X3DHKeyBundle{}, false
	}
	otp := c.Bundle.OtpSet[0]
	c.Bundle.OtpSet = c.Bundle.OtpSet[1:]
	// Return
	return x3dh_core.X3DHKeyBundle{
		IK:  ik,
		SPK: spk,
		OTP: otp,
	}, true
}

func (s *Server) SendMessage(clientID string, msg x3dh_core.InitialMessage) bool {
	c, ok := s.clients[clientID]
	if !ok {
		return false
	}
	c.Queue = append(c.Queue, msg)
	return true
}

func (s *Server) GetMessage(clientID string) (x3dh_core.InitialMessage, bool) {
	c, ok := s.clients[clientID]
	if !ok {
		return x3dh_core.InitialMessage{}, false
	}
	if len(c.Queue) == 0 {
		return x3dh_core.InitialMessage{}, false
	}
	msg := c.Queue[0]
	c.Queue = c.Queue[1:]
	return msg, true
}
