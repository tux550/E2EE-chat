package main

// Triple Diffie-Hellman Exchange (3DHX)
import (
	"fmt"
	"x3dh_client"
	"x3dh_server"
)

func main() {
	// SERVER
	s, err := x3dh_server.InitServer()
	if err != nil {
		fmt.Println(err)
		return
	}

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
	s.AddClient("Alice", bundle1)

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
	s.AddClient("Bob", bundle2)

	// Send message
	// Alice retrieves Bob's key bundle
	kb = s.GetClientBundle("Bob")
	// Alice sends message to Bob
	fmt.Println("Sending message...")
	msgData, err := c2.SendMessage(kb, []byte("Hello, World!"))
	if err != nil {
		fmt.Println(err)
		return
	}
	s.SendMessage("Bob", msgData)
	// Bob recieves message
	fmt.Println("Recieving message...")
	msg, err := c1.RecieveMessage(msgData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(msg))
	fmt.Println("Message recieved successfully!")
}
