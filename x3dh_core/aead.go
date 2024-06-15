package x3dh_core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 16
	keySize    = 32 // AES-256
	iterations = 100000
)

// Derive a key from a given secret using PBKDF2
func deriveKey(secret, salt []byte) []byte {
	return pbkdf2.Key(secret, salt, iterations, keySize, sha256.New)
}

// Encrypts the plaintext using AES-GCM with the given secret and associated data
func EncryptAEAD(secret, plaintext, associatedData []byte) (salt, nonce, ciphertext []byte, err error) {
	// Generate a random salt
	salt = make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, nil, err
	}

	// Derive a key from the secret and salt
	key := deriveKey(secret, salt)

	// Create AES block cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create GCM AEAD cipher
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, err
	}

	// Generate a random nonce
	nonce = make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, err
	}

	// Encrypt the plaintext
	ciphertext = aead.Seal(nil, nonce, plaintext, associatedData)
	return salt, nonce, ciphertext, nil
}

// Decrypts the ciphertext using AES-GCM with the given secret, salt, nonce, and associated data
func DecryptAEAD(secret, salt, nonce, ciphertext, associatedData []byte) (plaintext []byte, err error) {
	// Derive the key from the secret and salt
	key := deriveKey(secret, salt)

	// Create AES block cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM AEAD cipher
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt the ciphertext
	plaintext, err = aead.Open(nil, nonce, ciphertext, associatedData)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
