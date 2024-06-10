package main

// Eliptic curve diffie hellman key exchange example
import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
)

func generateKeys() (*ecdh.PrivateKey, error) {
	// Curve
	curve := ecdh.X25519()
	// Generate private key
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Return both keys
	return privateKey, nil
}

func main() {
	// Generate private key 1
	privateKey1, err := generateKeys()
	if err != nil {
		return
	}
	// Generate private key 2
	privateKey2, err := generateKeys()
	if err != nil {
		return
	}
	// Generate shared key
	sharedKey1, err := privateKey1.ECDH(privateKey2.PublicKey())
	if err != nil {
		fmt.Println(err)
		return
	}
	sharedKey2, err := privateKey2.ECDH(privateKey1.PublicKey())
	if err != nil {
		fmt.Println(err)
		return
	}
	// Check if both keys are the same
	if bytes.Equal(sharedKey1, sharedKey2) {
		fmt.Println("Keys are equal")
	} else {
		fmt.Println("Keys are not equal")
	}
	// Print shared key
	fmt.Println(sharedKey1)
	fmt.Println(sharedKey2)
}
