package x3dh

import (
	"errors"
	"fmt"

	"go.step.sm/crypto/x25519"
)

type X3DHClient struct {
	// Private Bundle
	KeyBundle *X3DHKeyBundle
}

type InitialMessage struct {
	// Identity Key
	IdentityKey x25519.PublicKey
	// Ephemeral Key
	EphemeralKey x25519.PublicKey
	// One Time Pre Key ID
	OneTimePreKeyID int
	// AEAD
	Ciphertext []byte
	AD         []byte
	// Nonce
	Nonce []byte
	Salt  []byte
}

func NewClient() *X3DHClient {
	return &X3DHClient{
		KeyBundle: &X3DHKeyBundle{},
	}
}

func InitClient() (*X3DHClient, error) {
	// Create a new client
	c := NewClient()
	// Generate key bundle
	kb, err := GenereateKeyBundle()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	c.KeyBundle = kb
	// Return the client
	return c, nil
}

func (c *X3DHClient) GetPublicKeyBundle() *X3DHPublicKeyBundle {
	return c.KeyBundle.PublicBundle()
}

func (c *X3DHClient) SendMessage(pkb *X3DHPublicKeyBundle, msg []byte) (*InitialMessage, error) {
	// Validate the key bundle
	valid := pkb.Validate()
	if !valid {
		return nil, errors.New("invalid key bundle")
	}
	// Generate Ephermal Key
	ephemeralKey, err := GenerateKeyPairX25519()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// TODO: Select a one time pre key
	otp_id := 0
	// Generate shared secret
	// kb.signedPreKey to []byte
	dh1, err := c.KeyBundle.ik.IdentityKey.SharedKey(pkb.spk.SignedPreKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dh2, err := ephemeralKey.SharedKey(pkb.ik.IdentityKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dh3, err := ephemeralKey.SharedKey(pkb.otp_set[otp_id].OneTimePreKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Concatenate the shared secrets
	sharedSecret := []byte{}
	sharedSecret = append(sharedSecret, dh1[:]...)
	sharedSecret = append(sharedSecret, dh2[:]...)
	sharedSecret = append(sharedSecret, dh3[:]...)
	// Build AD
	ad := []byte{}
	ad = append(ad, c.KeyBundle.ik.IdentityKey.PublicKey[:]...)
	ad = append(ad, pkb.ik.IdentityKey[:]...)
	// TODO: Encrypt the message with the shared secret using AEAD schema (msg encrypted + ad)
	salt, nonce, ciphertext, err := encryptAEAD(
		sharedSecret,
		msg,
		ad,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Generate a new initial message
	im := &InitialMessage{
		IdentityKey:     c.KeyBundle.ik.IdentityKey.PublicKey,
		EphemeralKey:    ephemeralKey.PublicKey,
		OneTimePreKeyID: otp_id,
		Ciphertext:      ciphertext,
		AD:              ad,
		Nonce:           nonce,
		Salt:            salt,
	}
	// Return the initial message
	return im, nil
}

func (c *X3DHClient) RecieveMessage(im *InitialMessage) ([]byte, error) {
	// Generate shared secret
	// kb.signedPreKey to []byte
	dh1, err := c.KeyBundle.spk.SignedPreKey.PrivateKey.SharedKey(im.IdentityKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dh2, err := c.KeyBundle.ik.IdentityKey.PrivateKey.SharedKey(im.EphemeralKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dh3, err := c.KeyBundle.otp_set[im.OneTimePreKeyID].OneTimePreKey.PrivateKey.SharedKey(im.EphemeralKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Concatenate the shared secrets
	sharedSecret := []byte{}
	sharedSecret = append(sharedSecret, dh1[:]...)
	sharedSecret = append(sharedSecret, dh2[:]...)
	sharedSecret = append(sharedSecret, dh3[:]...)
	// Build AD
	ad := []byte{}
	ad = append(ad, im.IdentityKey[:]...)
	ad = append(ad, c.KeyBundle.ik.IdentityKey.PublicKey[:]...)
	// Decrypt the message with the shared secret using AEAD schema (msg encrypted + ad)
	plaintext, err := decryptAEAD(
		sharedSecret,
		im.Salt,
		im.Nonce,
		im.Ciphertext,
		ad,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Return the plaintext
	return plaintext, nil
}
