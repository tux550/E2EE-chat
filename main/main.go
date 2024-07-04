package main

// Triple Diffie-Hellman Exchange (3DHX)
import (
	"encoding/json"
	"fmt"

	x3dh_client "tux.tech/x3dh/client"
	x3dh_core "tux.tech/x3dh/core"
	x3dh_server "tux.tech/x3dh/server"
)

func menu() {
	fmt.Println("===E2EE CHAT===")
	fmt.Println("1. Create account")
	fmt.Println("2. Log In")
	fmt.Println("Option >")

	var option int
	option, err := fmt.Scanln(&option);
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
    }
	switch option {
	case 1:
		createAccount()
	case 2:
		logIn()
	case 3:
		fmt.Println("~Goodbye!")
		return
	default:
		return
	}
}

func main() {
	// for {
	// 	menu()
	// }

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
	msgJSON, err := json.Marshal(*msgData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(msgJSON))
	msgUnmarshalled := x3dh_core.InitialMessage{}
	err = json.Unmarshal(msgJSON, &msgUnmarshalled)
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
