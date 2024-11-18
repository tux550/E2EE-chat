package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	X3DHCore "tux.tech/x3dh/core"
)

// ENV VARIABLES:
// - MONGO_URI
// - MONGO_DB
// - MONGO_CLIENT_COL
// - MONGO_MESSAGE_COL

type DocDBManager struct {
	client     *mongo.Client
	db         *mongo.Database
	clientCol  *mongo.Collection
	messageCol *mongo.Collection
}

func init() {
	// Validate environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
		panic("MONGO_URI environment variable is not set")
	}
	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		log.Fatal("MONGO_DB environment variable is not set")
		panic("MONGO_DB environment variable is not set")
	}
	mongoClientCol := os.Getenv("MONGO_CLIENT_COL")
	if mongoClientCol == "" {
		log.Fatal("MONGO_CLIENT_COL environment variable is not set")
		panic("MONGO_CLIENT_COL environment variable is not set")
	}
	mongoMessageCol := os.Getenv("MONGO_MESSAGE_COL")
	if mongoMessageCol == "" {
		log.Fatal("MONGO_MESSAGE_COL environment variable is not set")
		panic("MONGO_MESSAGE_COL environment variable is not set")
	}
}

func NewDocDBManager() *DocDBManager {
	// MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")
	mongoClientCol := os.Getenv("MONGO_CLIENT_COL")
	mongoMessageCol := os.Getenv("MONGO_MESSAGE_COL")

	// Connect to DocumentDB (MongoDB-compatible)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to DocumentDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping DocumentDB: %v", err)
	}

	log.Println("Connected to DocumentDB")

	// db
	db := client.Database(mongoDB)
	// clientCol
	clientCol := db.Collection(mongoClientCol)
	messageCol := db.Collection(mongoMessageCol)

	return &DocDBManager{
		client:     client,
		db:         db,
		clientCol:  clientCol,
		messageCol: messageCol,
	}
}

// === CLIENT COLLECTION ===

func (m *DocDBManager) ClientUpsert(clientID string, data ClientData) error {
	_, err := m.clientCol.UpdateOne(
		context.TODO(),
		map[string]string{"clientID": clientID},
		map[string]ClientData{"$set": data},
		options.Update().SetUpsert(true),
	)
	return err
}

func (m *DocDBManager) ClientExists(clientID string) (bool, error) {
	count, err := m.clientCol.CountDocuments(
		context.TODO(),
		map[string]string{"clientID": clientID},
	)
	return count > 0, err
}

func (m *DocDBManager) ClientFindOne(clientID string, data ClientData) error {
	return m.clientCol.FindOne(
		context.TODO(),
		map[string]string{"clientID": clientID},
	).Decode(data)
}

func (m *DocDBManager) ClientOTPCount(clientID string) (int, error) {
	var data ClientData
	err := m.clientCol.FindOne(
		context.TODO(),
		map[string]string{"clientID": clientID},
	).Decode(&data)
	if err != nil {
		return 0, err
	}
	return len(data.Bundle.OtpSet), nil
}

func (m *DocDBManager) ClienOTPAppend(clientID string, otps []X3DHCore.X3DHPublicOTP) error {
	// Use findOneAndUpdate to append OTPs atomically
	prev_doc := m.clientCol.FindOneAndUpdate(
		context.TODO(),
		map[string]string{"clientID": clientID}, // Filter
		bson.M{
			"$push": bson.M{"bundle.otpset": bson.M{"$each": otps}}, // Append new OTPs
		},
	)
	if prev_doc.Err() != nil {
		log.Printf("Error appending OTPs to client %s: %v", clientID, prev_doc.Err())
		return prev_doc.Err()
	}
	return nil
}

func (m *DocDBManager) ClienOTPPop(clientID string, count int) (*ClientData, error) {
	// Use findOneAndUpdate to pop OTPs atomically
	prev_doc := m.clientCol.FindOneAndUpdate(
		context.TODO(),
		map[string]string{"clientID": clientID}, // Filter
		bson.M{
			"$pop": bson.M{"bundle.otpset": -count}, // Pop OTPs
		},
		// Set options
		//options.FindOneAndUpdate().SetBypassDocumentValidation(true),
	)
	if prev_doc.Err() != nil {
		log.Printf("Error popping OTPs from client %s: %v", clientID, prev_doc.Err())
		return nil, prev_doc.Err()
	}
	// Decode the previous document
	var prev_data ClientData
	err := prev_doc.Decode(&prev_data)
	if err != nil {
		log.Printf("Error decoding previous document: %v", err)
		return nil, err
	}
	return &prev_data, nil
}

// === MESSAGE COLLECTION ===
func (m *DocDBManager) ClientMessageAppend(
	recipientID string,
	message MessageData) error {
	// Use findOneAndUpdate to append Message atomically
	prev_doc := m.clientCol.FindOneAndUpdate(
		context.TODO(),
		map[string]string{"clientID": recipientID}, // Filter
		bson.M{
			"$push": bson.M{"queue": message}, // Append new message
		},
	)
	// Check for errors
	if prev_doc.Err() != nil {
		log.Printf("Error appending message to client %s: %v", message.SenderID, prev_doc.Err())
		return prev_doc.Err()
	}
	return nil
}

func (m *DocDBManager) ClientMessagePull(clientID string) ([]MessageData, error) {
	// Use findOneAndUpdate to pull all messages atomically
	prev_doc := m.clientCol.FindOneAndUpdate(
		context.TODO(),
		map[string]string{"clientID": clientID}, // Filter
		bson.M{
			"$set": bson.M{"queue": []MessageData{}}, // Clear the queue
		},
	)
	// Check for errors
	if prev_doc.Err() != nil {
		log.Printf("Error pulling messages from client %s: %v", clientID, prev_doc.Err())
		return nil, prev_doc.Err()
	}
	// Decode the previous document
	var prev_data ClientData
	err := prev_doc.Decode(&prev_data)
	if err != nil {
		log.Printf("Error decoding previous document: %v", err)
		return nil, err
	}
	return prev_data.Queue, nil
}

func (m *DocDBManager) ClientMessagePop(clientID string, count int) ([]MessageData, error) {
	// Use findOneAndUpdate to pop messages atomically
	prev_doc := m.clientCol.FindOneAndUpdate(
		context.TODO(),
		map[string]string{"clientID": clientID}, // Filter
		bson.M{
			"$pop": bson.M{"queue": -count}, // Pop messages
		},
	)
	// Check for errors
	if prev_doc.Err() != nil {
		log.Printf("Error popping messages from client %s: %v", clientID, prev_doc.Err())
		return nil, prev_doc.Err()
	}
	// Decode the previous document
	var prev_data ClientData
	err := prev_doc.Decode(&prev_data)
	if err != nil {
		log.Printf("Error decoding previous document: %v", err)
		return nil, err
	}
	return prev_data.Queue, nil
}
