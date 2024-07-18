package main

import (
	"fmt"
	"os"

	x3dh_client "tux.tech/x3dh/client"
)

var url = "ws://localhost:8765/ws"
var username = "Bob"

var secrets_filename = ".secrets.json"

func GetMyClient() (*x3dh_client.X3DHClient, error) {
	// Check if secrets file exists
	_, err := os.Stat(secrets_filename)
	if os.IsNotExist(err) {
		// Create new secrets client
		fmt.Println("Initialized new client")
		client, err := x3dh_client.InitClient()
		if err != nil {
			return nil, err
		}
		return client, nil
	} else {
		// Load secrets from file
		fmt.Println("Loaded existing client")
		client, err := x3dh_client.LoadClient(secrets_filename)
		if err != nil {
			return nil, err
		}
		return client, nil
	}
}

func SaveMyClient(client *x3dh_client.X3DHClient) error {
	// Save secrets to file
	fmt.Println("Saved client")
	err := client.SaveClient(secrets_filename)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	a, err := GetMyClient()

	if err != nil {
		fmt.Println("Could not get client:", err)
		return
	}

	a.DebugPrint()

	err = SaveMyClient(a)

	if err != nil {
		fmt.Println("Could not save client:", err)
		return
	}

	// Print Client

	/*
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
	*/
}
