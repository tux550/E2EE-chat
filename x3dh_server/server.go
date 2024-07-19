package x3dh_server

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	db        *mongo.Database
	clientCol *mongo.Collection
	//clients map[string]*ClientData
}

func NewServer() *Server {
	// Create connection to mongo
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	db := client.Database("x3dh")
	clientCol := db.Collection("clients")

	return &Server{
		db:        db,
		clientCol: clientCol,

		//clients: make(map[string]*ClientData),
	}
}

func NewClientData(bundle X3DHCore.X3DHClientBundle) *ClientData {
	return &ClientData{
		Bundle: bundle,
		Queue:  make([]MessageData, 0),
	}
}

/*
func (s *Server) RegisterClient(clientID string, bundle X3DHCore.X3DHClientBundle) {
	s.clients[clientID] = NewClientData(bundle)
}*/

func (s *Server) RegisterClient(clientID string, bundle X3DHCore.X3DHClientBundle) error {
	data := NewClientData(bundle)
	_, err := s.clientCol.UpdateOne(
		context.TODO(),
		bson.M{"clientID": clientID},
		bson.M{"$set": data},
		options.Update().SetUpsert(true),
	)
	return err
}

/*func (s *Server) IsClientRegistered(clientID string) bool {
	_, ok := s.clients[clientID]
	return ok
}*/

func (s *Server) IsClientRegistered(clientID string) (bool, error) {
	count, err := s.clientCol.CountDocuments(
		context.TODO(),
		bson.M{"clientID": clientID},
	)
	return count > 0, err
}

/*func (s *Server) GetRemainingOTPCount(clientID string) int {
	c, ok := s.clients[clientID]
	if !ok {
		return 0
	}
	return len(c.Bundle.OtpSet)
}*/

func (s *Server) GetRemainingOTPCount(clientID string) (int, error) {
	var clientData ClientData
	err := s.clientCol.FindOne(
		context.TODO(),
		bson.M{"clientID": clientID},
	).Decode(&clientData)
	if err != nil {
		return 0, err
	}
	return len(clientData.Bundle.OtpSet), nil
}

/*func (s *Server) ExpandOTPSet(clientID string, otps []X3DHCore.X3DHPublicOTP) {
	c, ok := s.clients[clientID]
	if !ok {
		return
	}
	c.Bundle.OtpSet = append(c.Bundle.OtpSet, otps...)
}*/

func (s *Server) ExpandOTPSet(clientID string, otps []X3DHCore.X3DHPublicOTP) error {
	var clientData ClientData
	err := s.clientCol.FindOne(
		context.TODO(),
		bson.M{"clientID": clientID},
	).Decode(&clientData)
	if err != nil {
		return err
	}

	clientData.Bundle.OtpSet = append(clientData.Bundle.OtpSet, otps...)
	_, err = s.clientCol.UpdateOne(
		context.TODO(),
		bson.M{"clientID": clientID},
		bson.M{"$set": bson.M{"bundle.otpSet": clientData.Bundle.OtpSet}},
	)
	return err
}

/*func (s *Server) GetClientBundle(clientID string) (X3DHCore.X3DHKeyBundle, bool) {
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
}*/

func (s *Server) GetClientBundle(clientID string) (X3DHCore.X3DHKeyBundle, bool, error) {
	var clientData ClientData
	err := s.clientCol.FindOne(
		context.TODO(),
		bson.M{"clientID": clientID},
	).Decode(&clientData)
	if err != nil {
		return X3DHCore.X3DHKeyBundle{}, false, err
	}

	if len(clientData.Bundle.OtpSet) == 0 {
		return X3DHCore.X3DHKeyBundle{}, false, nil
	}

	otp := clientData.Bundle.OtpSet[0]
	clientData.Bundle.OtpSet = clientData.Bundle.OtpSet[1:]

	_, err = s.clientCol.UpdateOne(
		context.TODO(),
		bson.M{"clientID": clientID},
		bson.M{"$set": bson.M{"bundle.otpSet": clientData.Bundle.OtpSet}},
	)
	if err != nil {
		return X3DHCore.X3DHKeyBundle{}, false, err
	}

	return X3DHCore.X3DHKeyBundle{
		IK:  clientData.Bundle.IK,
		SPK: clientData.Bundle.SPK,
		OTP: otp,
	}, true, nil
}

/*func (s *Server) SendMessage(recipientID string, senderID string, msg X3DHCore.InitialMessage) bool {
	c, ok := s.clients[recipientID]
	if !ok {
		return false
	}
	c.Queue = append(c.Queue, MessageData{
		SenderID: senderID,
		Message:  msg,
	})
	return true
}*/

func (s *Server) SendMessage(recipientID string, senderID string, msg X3DHCore.InitialMessage) bool {
	_, err := s.clientCol.UpdateOne(
		context.TODO(),
		bson.M{"clientID": recipientID},
		bson.M{"$push": bson.M{"queue": MessageData{
			SenderID: senderID,
			Message:  msg,
		}}},
	)
	return err == nil
}

/*func (s *Server) GetMessage(clientID string) (MessageData, bool) {
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
}*/

func (s *Server) GetMessage(clientID string) (MessageData, bool, error) {
	var clientData ClientData
	err := s.clientCol.FindOne(
		context.TODO(),
		bson.M{"clientID": clientID},
	).Decode(&clientData)
	if err != nil {
		return MessageData{}, false, err
	}

	if len(clientData.Queue) == 0 {
		return MessageData{}, false, nil
	}

	msg := clientData.Queue[0]
	clientData.Queue = clientData.Queue[1:]

	_, err = s.clientCol.UpdateOne(
		context.TODO(),
		bson.M{"clientID": clientID},
		bson.M{"$set": bson.M{"queue": clientData.Queue}},
	)
	if err != nil {
		return MessageData{}, false, err
	}

	return msg, true, nil
}
