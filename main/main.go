package main

// Triple Diffie-Hellman Exchange (3DHX)
import (
	"fmt"
	"x3dh"
)

func main() {
	c1, err := x3dh.InitClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	bundle := c1.GetPublicKeyBundle()

	c2, err := x3dh.InitClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	msgData, err := c2.SendMessage(bundle, []byte("Hello, World!"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sending message...")
	fmt.Println("Recieving message...")
	msg, err := c1.RecieveMessage(msgData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(msg))
	fmt.Println("Message recieved successfully!")
}
