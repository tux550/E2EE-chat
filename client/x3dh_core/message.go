package x3dh_core

import "go.step.sm/crypto/x25519"

type InitialMessage struct {
	// Identity Key
	IdentityKey x25519.PublicKey `json:"identity_key"`
	// Ephemeral Key
	EphemeralKey x25519.PublicKey `json:"ephemeral_key"`
	// One Time Pre Key ID
	OneTimePreKeyID int `json:"one_time_pre_key_id"`
	// AEAD
	Ciphertext []byte `json:"ciphertext"`
	AD         []byte `json:"ad"`
	// Nonce
	Nonce []byte `json:"nonce"`
	Salt  []byte `json:"salt"`
}
