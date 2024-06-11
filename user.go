package main

// import go.step.sm/crypto

import (
	"crypto"
	"crypto/rand"
	"fmt"

	"go.step.sm/crypto/x25519"
)

type KeyPair struct {
	// Public key
	PublicKey x25519.PublicKey
	// Private key
	PrivateKey x25519.PrivateKey
}

type KeyBundle struct {
	// Identity Key
	IdentityKey x25519.PublicKey
	// Signed Pre Key
	SignedPreKey x25519.PublicKey
	// One Time Pre Key
	OneTimePreKeys []x25519.PublicKey
	// One Time Pre Key IDs
	OneTimePreKeyIDs []int
	// Signed Pre Key Signature
	SignedPreKeySignature []byte
}

type User struct {
	// Identity Key
	IdentityKey *KeyPair
	// Signed Pre Key
	SignedPreKey *KeyPair
	// One Time Pre Key Dictionary
	OneTimePreKey map[int]*KeyPair
	// One Time Pre Key counter
	OneTimePreKeyCounter int
	// Signed Pre Key Signature
	SignedPreKeySignature []byte
}

func NewUser() *User {
	return &User{
		OneTimePreKey: make(map[int]*KeyPair),
	}
}

func InitUser() (*User, error) {
	// Create a new user
	u := NewUser()
	// Generate a new identity key
	identityKey, err := newKeyPair()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	u.IdentityKey = identityKey
	// Generate a new signed pre key
	signedPreKey, err := newKeyPair()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	u.SignedPreKey = signedPreKey
	// Generate a new one time pre key
	u.OneTimePreKeyCounter = 0
	for i := 0; i < 5; i++ {
		err := u.GenerateOneTimePreKey()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	// Return the user
	return u, nil
}

func (u *User) GenerateOneTimePreKey() error {
	// Generate a new one time pre key
	oneTimePreKey, err := newKeyPair()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Add the one time pre key to the dictionary
	u.OneTimePreKey[u.OneTimePreKeyCounter] = oneTimePreKey
	// Increment the counter
	u.OneTimePreKeyCounter++
	// Return
	return nil
}

func (u *User) GenerateKeyBundle() *KeyBundle {
	// Create signature
	signature, err := signPreKey(u.IdentityKey.PrivateKey, u.SignedPreKey.PublicKey)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Create a new key bundle
	kb := &KeyBundle{
		IdentityKey:           u.IdentityKey.PublicKey,
		SignedPreKey:          u.SignedPreKey.PublicKey,
		OneTimePreKeys:        make([]x25519.PublicKey, 0),
		OneTimePreKeyIDs:      make([]int, 0),
		SignedPreKeySignature: signature,
	}
	// Add the one time pre keys to the key bundle
	for k, v := range u.OneTimePreKey {
		kb.OneTimePreKeys = append(kb.OneTimePreKeys, v.PublicKey)
		kb.OneTimePreKeyIDs = append(kb.OneTimePreKeyIDs, k)
	}
	// Return the key bundle
	return kb
}

func newKeyPair() (*KeyPair, error) {
	// Generate private key
	publicKey, privateKey, err := x25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Return
	return &KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func signPreKey(PrivateIdentityKey x25519.PrivateKey, PublicPreKey x25519.PublicKey) ([]byte, error) {
	// Sign the public key
	signature, err := PrivateIdentityKey.Sign(rand.Reader, PublicPreKey, crypto.Hash(0)) // Review: https://pkg.go.dev/go.step.sm/crypto/x25519
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Return
	return signature, nil
}