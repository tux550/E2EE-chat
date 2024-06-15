package x3dh_core

import "go.step.sm/crypto/x25519"

type X3DHFullIK struct {
	// Identity Key
	IdentityKey KeyPairX25519
}

type X3DHPublicIK struct {
	// Identity Key
	IdentityKey x25519.PublicKey
}

func (ik *X3DHFullIK) PublicIK() *X3DHPublicIK {
	return &X3DHPublicIK{
		IdentityKey: ik.IdentityKey.PublicKey,
	}
}

func GenerateFullIK() (*X3DHFullIK, error) {
	// Generate private key
	kp, err := GenerateKeyPairX25519()
	if err != nil {
		return nil, err
	}
	// Return
	return &X3DHFullIK{
		IdentityKey: *kp,
	}, nil
}
