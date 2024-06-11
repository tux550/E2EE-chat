package main

// Triple Diffie-Hellman Exchange (3DHX)
import (
	"fmt"
)

func main() {
	user, err := InitUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	bundle := user.GenerateKeyBundle()
	fmt.Println("Bundle:", bundle)
}
