package main

// Triple Diffie-Hellman Exchange (3DHX)
import (
	"fmt"
)

func main() {
	user1, err := InitUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	bundle := user1.GetKeyBundle()

	user2, err := InitUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	msgData, err := user2.SendMessage(bundle, []byte("Hello, World!"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(msgData)
}
