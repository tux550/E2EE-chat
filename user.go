package main

// import go.step.sm/crypto

import (
	"crypto"
	"crypto/rand"
	"errors"
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

type InitialMessage struct {
	// Identity Key
	IdentityKey x25519.PublicKey
	// Ephemeral Key
	EphemeralKey x25519.PublicKey
	// One Time Pre Key ID
	OneTimePreKeyID int
	// AEAD
	AEAD []byte
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

func (u *User) GetKeyBundle() *KeyBundle {
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

func (kb *KeyBundle) Validate() bool {
	// Validate the signed pre key
	valid := x25519.Verify(kb.IdentityKey, kb.SignedPreKey, kb.SignedPreKeySignature)
	return valid
}

func (u *User) SendMessage(kb *KeyBundle, msg []byte) (*InitialMessage, error) {
	// Validate the key bundle
	valid := kb.Validate()
	if !valid {
		fmt.Println("Invalid key bundle")
		return nil, errors.New("Invalid key bundle")
	}
	fmt.Println("Valid key bundle")
	// Generate Ephermal Key
	ephemeralKey, err := newKeyPair()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// TODO: Select a one time pre key
	otp_id := 0
	// Generate shared secret
	// kb.signedPreKey to []byte
	dh1, err := u.IdentityKey.PrivateKey.SharedKey(kb.SignedPreKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dh2, err := ephemeralKey.PrivateKey.SharedKey(kb.IdentityKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dh3, err := ephemeralKey.PrivateKey.SharedKey(kb.OneTimePreKeys[otp_id])
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Concatenate the shared secrets
	sharedSecret := []byte{}
	sharedSecret = append(sharedSecret, dh1[:]...)
	sharedSecret = append(sharedSecret, dh2[:]...)
	sharedSecret = append(sharedSecret, dh3[:]...)
	// TODO: Derive a key from the shared secret

	// Build AD
	ad := []byte{}
	ad = append(ad, u.IdentityKey.PublicKey[:]...)
	ad = append(ad, kb.IdentityKey[:]...)
	// TODO: Encrypt the message with the shared secret using AEAD schema (msg encrypted + ad)

	// Generate a new initial message
	im := &InitialMessage{
		IdentityKey:     u.IdentityKey.PublicKey,
		EphemeralKey:    ephemeralKey.PublicKey,
		OneTimePreKeyID: kb.OneTimePreKeyIDs[0],
		AEAD:            []byte{}, // TODO: Encrypted message
	}
	// Return the initial message
	return im, nil
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
