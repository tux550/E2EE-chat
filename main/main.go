package main

// Triple Diffie-Hellman Exchange (3DHX)
import (
	"fmt"
	"x3dh_client"
	"x3dh_server"
)

func main() {
	// SERVER
	s := x3dh_server.NewServer()

	// CLIENTS
	c1, err := x3dh_client.InitClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	bundle1, err := c1.GetServerInitBundle()
	if err != nil {
		fmt.Println(err)
		return
	}
	s.RegisterClient("Alice", *bundle1)
	c2, err := x3dh_client.InitClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	bundle2, err := c2.GetServerInitBundle()
	if err != nil {
		fmt.Println(err)
		return
	}
	s.RegisterClient("Bob", *bundle2)
	// Send message
	// Alice retrieves Bob's key bundle
	kb, ok := s.GetClientBundle("Bob")
	if !ok {
		fmt.Println("Bob's key bundle not found")
		return
	}
	// Alice sends message to Bob
	fmt.Println("Sending message...")
	msgData, err := c1.BuildMessage(&kb, []byte("Hello, World!"))
	if err != nil {
		fmt.Println(err)
		return
	}
	ok = s.SendMessage("Bob", *msgData)
	if !ok {
		fmt.Println("Message not sent")
		return
	}
	bobMessage, ok := s.GetMessage("Bob")
	if !ok {
		fmt.Println("Bob's message not found")
		return
	}
	// Bob recieves message
	fmt.Println("Recieving message...")
	msg, err := c2.RecieveMessage(&bobMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(msg))
	fmt.Println("Message recieved successfully!")
}
