package x3dh

import (
	"crypto"
	"crypto/rand"
	"fmt"

	"go.step.sm/crypto/x25519"
)

type KeyPairX25519 struct {
	// Public key
	PublicKey x25519.PublicKey
	// Private key
	PrivateKey x25519.PrivateKey
}

func GenerateKeyPairX25519() (*KeyPairX25519, error) {
	// Generate private key
	publicKey, privateKey, err := x25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Return
	return &KeyPairX25519{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func (kp *KeyPairX25519) Sign(value []byte) ([]byte, error) {
	// Sign value
	signature, err := kp.PrivateKey.Sign(rand.Reader, value, crypto.Hash(0)) // Review: https://pkg.go.dev/go.step.sm/crypto/x25519
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Return
	return signature, nil
}

func (kp *KeyPairX25519) SharedKey(publicKey x25519.PublicKey) ([]byte, error) {
	// Generate shared key
	sharedKey, err := kp.PrivateKey.SharedKey(publicKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Return
	return sharedKey, nil
}
