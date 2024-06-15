package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var url = "ws://localhost:8765/ws"
var username = "Bob"

func main() {
	// Add username to header
	header := http.Header{}
	header.Add("User", username)

	c, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Handle incoming messages
	go func() {
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Could not read message:", err)
				return
			}
			if mt == websocket.BinaryMessage {
				fmt.Println("Received message:", message)
			}
		}
	}()

	// Send message
	message := []byte("Hello, world!")
	err = c.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		fmt.Println("Could not send message:", err)
		return
	}
	fmt.Println("Sent message:", message)

	// Graceful close
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println("Could not send close message:", err)
		return
	}
	fmt.Println("Sent close message")
}
