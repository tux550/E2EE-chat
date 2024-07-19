package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	x3dh_server "tux.tech/x3dh/server"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsServer struct {
	dbClient   *mongo.Client
	X3DHServer x3dh_server.Server
	clients    map[*WsClient]bool
	mu         sync.Mutex
}

func NewWsServer() *WsServer {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	return &WsServer{
		clients:    make(map[*WsClient]bool),
		X3DHServer: *x3dh_server.NewServer(),
		dbClient:   client,
	}
}

func (server *WsServer) SetClient(client *WsClient) {
	server.mu.Lock()
	server.clients[client] = true
	server.mu.Unlock()
}

func (server *WsServer) UnsetClient(client *WsClient) {
	server.mu.Lock()
	delete(server.clients, client)
	server.mu.Unlock()
}

func (server *WsServer) SendNotificationToUser(user string, message []byte) {
	server.mu.Lock()
	for client := range server.clients {
		if client.username == user {
			client.send <- message
		}
	}
	server.mu.Unlock()
}

func (server *WsServer) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (server *WsServer) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (server *WsServer) authenticateUser(username, password string) bool {
	collection := server.dbClient.Database("yourdatabase").Collection("users")

	// Find the user in the database
	var result struct {
		Username string `bson:"username"`
		Password string `bson:"password"`
	}
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User does not exist, create a new user
			hashedPassword, err := server.hashPassword(password)
			if err != nil {
				fmt.Println("Error hashing password:", err)
				return false
			}

			newUser := bson.M{"username": username, "password": hashedPassword}
			_, err = collection.InsertOne(context.Background(), newUser)
			if err != nil {
				fmt.Println("Error creating new user:", err)
				return false
			}
			return true
		}
		fmt.Println("Error finding user:", err)
		return false
	}

	// Check if the password matches the hashed password
	return server.checkPasswordHash(password, result.Password)
}

func (server *WsServer) connnect(w http.ResponseWriter, r *http.Request) {
	// Get user fom header
	user := r.Header.Get("User")
	password := r.Header.Get("Password")

	if user == "" || password == "" {
		http.Error(w, "No user", http.StatusBadRequest)
		return
	}

	// Authenticate user
	if !server.authenticateUser(user, password) {
		http.Error(w, "Invalid auth", http.StatusUnauthorized)
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	client := NewWsClient(user, server, conn)

	server.SetClient(client)

	fmt.Println("New connection from", user)
	go client.WritePump()
	go client.ReadPump()
}
