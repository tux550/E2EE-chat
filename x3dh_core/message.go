package x3dh_core

import "go.step.sm/crypto/x25519"

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
