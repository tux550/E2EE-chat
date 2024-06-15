package x3dh_client

import (
	"errors"
	"x3dh_core"
)

type X3DHClient struct {
	// Identity
	IdentityKey x3dh_core.X3DHFullIK
	// Signed Pre Key
	SignedPreKey x3dh_core.X3DHFullSPK
	// One Time Pre Keys
	OneTimePreKeys []x3dh_core.X3DHFullOTP
	// Counter
	otpCounter int
}

func NewClient() *X3DHClient {
	return &X3DHClient{
		otpCounter: 0,
	}
}

// TODO: Load Client from file
// func LoadClient() (*X3DHClient, error) {}

func InitClient() (*X3DHClient, error) {
	// Create a new client
	c := NewClient()
	// Identity Key
	ik, err := x3dh_core.GenerateFullIK()
	if err != nil {
		return nil, err
	}
	c.IdentityKey = *ik
	// Signed Pre Key
	spk, err := x3dh_core.GenerateFullSPK(c.IdentityKey.IdentityKey)
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
	otp, err := x3dh_core.GenerateFullOTP(c.otpCounter)
	if err != nil {
		return err
	}
	// Append the one time pre key
	c.OneTimePreKeys = append(c.OneTimePreKeys, *otp)
	// Increment the counter
	c.otpCounter++
	// Return
	return nil
}

func (c *X3DHClient) GetServerInitBundle() (*x3dh_core.X3DHClientBundle, error) {
	// Generate the server key bundle
	otp_set := make([]x3dh_core.X3DHPublicOTP, 0)
	for _, otp := range c.OneTimePreKeys {
		otp_set = append(otp_set, *otp.PublicOTP())
	}
	skb := &x3dh_core.X3DHClientBundle{
		IK:     *c.IdentityKey.PublicIK(),
		SPK:    *c.SignedPreKey.PublicSPK(),
		OtpSet: otp_set,
	}
	// Return the server key bundle
	return skb, nil
}

func (c *X3DHClient) BuildMessage(pkb *x3dh_core.X3DHKeyBundle, msg []byte) (*x3dh_core.InitialMessage, error) {
	// Validate the key bundle
	valid := pkb.Validate()
	if !valid {
		return nil, errors.New("invalid key bundle")
	}
	// Generate Ephermal Key
	ephemeralKey, err := x3dh_core.GenerateKeyPairX25519()
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
	salt, nonce, ciphertext, err := x3dh_core.EncryptAEAD(
		sharedSecret,
		msg,
		ad,
	)
	if err != nil {
		return nil, err
	}
	// Generate a new initial message
	im := &x3dh_core.InitialMessage{
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

func (c *X3DHClient) RecieveMessage(im *x3dh_core.InitialMessage) ([]byte, error) {
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
	plaintext, err := x3dh_core.DecryptAEAD(
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

/*
func (c *X3DHClient) SendMessage(pkb *X3DHClientKeyBundle, msg []byte) (*InitialMessage, error) {
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
	// Generate shared secret
	// kb.signedPreKey to []byte
	dh1, err := c.KeyBundle.ik.IdentityKey.SharedKey(pkb.SPK.SignedPreKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dh2, err := ephemeralKey.SharedKey(pkb.IK.IdentityKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dh3, err := ephemeralKey.SharedKey(pkb.OTP.OneTimePreKey)
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
	ad = append(ad, pkb.IK.IdentityKey[:]...)
	// Encrypt the message with the shared secret using AEAD schema (msg encrypted + ad)
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
		OneTimePreKeyID: pkb.OTP.OneTimePreKeyID,
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
*/
