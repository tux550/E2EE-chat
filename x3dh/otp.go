package x3dh

import "go.step.sm/crypto/x25519"

type X3DHFullOTP struct {
	// One Time Pre Key
	OneTimePreKey KeyPairX25519
	// One Time Pre Key ID
	OneTimePreKeyID int
}

type X3DHPublicOTP struct {
	// One Time Pre Key
	OneTimePreKey x25519.PublicKey
	// One Time Pre Key ID
	OneTimePreKeyID int
}

func (otp *X3DHFullOTP) PublicOTP() *X3DHPublicOTP {
	return &X3DHPublicOTP{
		OneTimePreKey:   otp.OneTimePreKey.PublicKey,
		OneTimePreKeyID: otp.OneTimePreKeyID,
	}
}

func GenerateFullOTP(id int) (*X3DHFullOTP, error) {
	// Generate private key
	kp, err := GenerateKeyPairX25519()
	if err != nil {
		return nil, err
	}
	// Return
	return &X3DHFullOTP{
		OneTimePreKey:   *kp,
		OneTimePreKeyID: id,
	}, nil
}
