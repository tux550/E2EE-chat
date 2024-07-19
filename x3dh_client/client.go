package x3dh_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	X3DHCore "tux.tech/x3dh/core"
)

type X3DHClient struct {
	// Username
	Username string `json:"username"`
	// Identity
	IdentityKey X3DHCore.X3DHFullIK `json:"identityKey"`
	// Signed Pre Key
	SignedPreKey X3DHCore.X3DHFullSPK `json:"signedPreKey"`
	// One Time Pre Keys
	OneTimePreKeys []X3DHCore.X3DHFullOTP `json:"oneTimePreKeys"`
	// Counter
	OTPCounter int `json:"otpCounter"`
}

func NewClient() *X3DHClient {
	return &X3DHClient{
		OTPCounter: 0,
	}
}

func (c *X3DHClient) DebugPrint() {
	// Print the client
	fmt.Println("=== Client ===")
	fmt.Println("Username: ", c.Username)
	fmt.Println("Identity Key: ", c.IdentityKey.IdentityKey.PublicKey)
	fmt.Println("Signed Pre Key: ", c.SignedPreKey.SignedPreKey.PublicKey)
	fmt.Println("One Time Pre Keys: ")
	for i, otp := range c.OneTimePreKeys {
		fmt.Println("One Time Pre Key ", i, ": ", otp.OneTimePreKey.PublicKey)
	}
	fmt.Println("OTP Counter: ", c.OTPCounter)
	fmt.Println("=== End ===")
}

// Save client
func (c *X3DHClient) SaveClient(target_filename string) error {
	// Marshal the client to JSON
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	// Write the data to the file
	err = os.WriteFile(target_filename, data, 0644)
	if err != nil {
		return err
	}
	// Return
	return nil
}

// Load client
func LoadClient(target_filename string) (*X3DHClient, error) {
	// Read the data from the file
	data, err := os.ReadFile(target_filename)
	if err != nil {
		return nil, err
	}
	// Unmarshal the data to a client
	c := NewClient()
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	// Return the client
	return c, nil
}

func InitClient(username string) (*X3DHClient, error) {
	// Create a new client
	c := NewClient()
	// Set the username
	c.Username = username
	// Identity Key
	ik, err := X3DHCore.GenerateFullIK()
	if err != nil {
		return nil, err
	}
	c.IdentityKey = *ik
	// Signed Pre Key
	spk, err := X3DHCore.GenerateFullSPK(c.IdentityKey.IdentityKey)
	if err != nil {
		return nil, err
	}
	c.SignedPreKey = *spk
	// One Time Pre Keys
	for i := 0; i < 5; i++ {
		err := c.generateOneTimePreKey()
		if err != nil {
			return nil, err
		}
	}
	// Return the client
	return c, nil
}

func (c *X3DHClient) generateOneTimePreKey() error {
	// Generate one time pre key
	otp, err := X3DHCore.GenerateFullOTP(c.OTPCounter)
	if err != nil {
		return err
	}
	// Append the one time pre key
	c.OneTimePreKeys = append(c.OneTimePreKeys, *otp)
	// Increment the counter
	c.OTPCounter += 1
	// Return
	return nil
}

func (c *X3DHClient) GetServerInitBundle() (*X3DHCore.X3DHClientBundle, error) {
	// Generate the server key bundle
	otp_set := make([]X3DHCore.X3DHPublicOTP, 0)
	for _, otp := range c.OneTimePreKeys {
		otp_set = append(otp_set, *otp.PublicOTP())
	}
	skb := &X3DHCore.X3DHClientBundle{
		IK:     *c.IdentityKey.PublicIK(),
		SPK:    *c.SignedPreKey.PublicSPK(),
		OtpSet: otp_set,
	}
	// Return the server key bundle
	return skb, nil
}

func (c *X3DHClient) BuildMessage(pkb *X3DHCore.X3DHKeyBundle, msg []byte) (*X3DHCore.InitialMessage, error) {
	// Validate the key bundle
	valid := pkb.Validate()
	if !valid {
		return nil, errors.New("invalid key bundle")
	}
	// Generate Ephermal Key
	ephemeralKey, err := X3DHCore.GenerateKeyPairX25519()
	if err != nil {
		return nil, err
	}
	// Generate shared secret
	// kb.signedPreKey to []byte
	dh1, err := c.IdentityKey.IdentityKey.SharedKey(pkb.SPK.SignedPreKey)
	if err != nil {
		return nil, err
	}
	dh2, err := ephemeralKey.SharedKey(pkb.IK.IdentityKey)
	if err != nil {
		return nil, err
	}
	dh3, err := ephemeralKey.SharedKey(pkb.OTP.OneTimePreKey)
	if err != nil {
		return nil, err
	}
	// Concatenate the shared secrets
	sharedSecret := []byte{}
	sharedSecret = append(sharedSecret, dh1[:]...)
	sharedSecret = append(sharedSecret, dh2[:]...)
	sharedSecret = append(sharedSecret, dh3[:]...)
	// Build AD
	ad := []byte{}
	ad = append(ad, c.IdentityKey.IdentityKey.PublicKey[:]...)
	ad = append(ad, pkb.IK.IdentityKey[:]...)
	// Encrypt the message with the shared secret using AEAD schema (msg encrypted + ad)
	salt, nonce, ciphertext, err := X3DHCore.EncryptAEAD(
		sharedSecret,
		msg,
		ad,
	)
	if err != nil {
		return nil, err
	}
	// Generate a new initial message
	im := &X3DHCore.InitialMessage{
		IdentityKey:     c.IdentityKey.IdentityKey.PublicKey,
		EphemeralKey:    ephemeralKey.PublicKey,
		OneTimePreKeyID: pkb.OTP.OneTimePreKeyID,
		Ciphertext:      ciphertext,
		AD:              ad,
		Nonce:           nonce,
		Salt:            salt,
	}
	// Return the initial message
	return im, nil
}

func (c *X3DHClient) RecieveMessage(im *X3DHCore.InitialMessage) ([]byte, error) {
	// Generate shared secret
	// kb.signedPreKey to []byte
	dh1, err := c.SignedPreKey.SignedPreKey.PrivateKey.SharedKey(im.IdentityKey)
	if err != nil {
		return nil, err
	}
	dh2, err := c.IdentityKey.IdentityKey.PrivateKey.SharedKey(im.EphemeralKey)
	if err != nil {
		return nil, err
	}
	dh3, err := c.OneTimePreKeys[im.OneTimePreKeyID].OneTimePreKey.PrivateKey.SharedKey(im.EphemeralKey)
	if err != nil {
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
	ad = append(ad, c.IdentityKey.IdentityKey.PublicKey[:]...)
	// Decrypt the message with the shared secret using AEAD schema (msg encrypted + ad)
	plaintext, err := X3DHCore.DecryptAEAD(
		sharedSecret,
		im.Salt,
		im.Nonce,
		im.Ciphertext,
		ad,
	)
	if err != nil {
		return nil, err
	}
	// Return the plaintext
	return plaintext, nil
}
