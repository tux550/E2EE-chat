package x3dh

import "go.step.sm/crypto/x25519"

type X3DHFullSPK struct {
	// Signed Pre Key
	SignedPreKey KeyPairX25519
	// Signed Pre Key Signature
	SignedPreKeySignature []byte
}

type X3DHPublicSPK struct {
	// Signed Pre Key
	SignedPreKey x25519.PublicKey
	// Signed Pre Key Signature
	SignedPreKeySignature []byte
}

func (spk *X3DHFullSPK) PublicSPK() *X3DHPublicSPK {
	return &X3DHPublicSPK{
		SignedPreKey:          spk.SignedPreKey.PublicKey,
		SignedPreKeySignature: spk.SignedPreKeySignature,
	}
}

func GenerateFullSPK(identityKey KeyPairX25519) (*X3DHFullSPK, error) {
	// Generate private key
	kp, err := GenerateKeyPairX25519()
	if err != nil {
		return nil, err
	}
	signature, err := identityKey.Sign(kp.PublicKey[:])
	if err != nil {
		return nil, err
	}
	// Return
	return &X3DHFullSPK{
		SignedPreKey:          *kp,
		SignedPreKeySignature: signature,
	}, nil
}
