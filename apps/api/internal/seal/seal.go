// Package seal encrypts the Matrix access tokens Prisme keeps, so a database
// dump alone hands out no session.
package seal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// Box seals and opens values with AES-256-GCM under one key.
type Box struct {
	aead cipher.AEAD
}

// New derives the key from a secret of any length.
func New(secret string) (*Box, error) {
	if len(secret) < 16 {
		return nil, errors.New("seal: the secret must be at least 16 characters")
	}
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Box{aead: aead}, nil
}

// Seal encrypts a value into a base64 string with a random nonce.
func (b *Box) Seal(plaintext string) string {
	nonce := make([]byte, b.aead.NonceSize())
	_, _ = rand.Read(nonce)
	return base64.RawStdEncoding.EncodeToString(b.aead.Seal(nonce, nonce, []byte(plaintext), nil))
}

// Open decrypts what Seal produced.
func (b *Box) Open(sealed string) (string, error) {
	data, err := base64.RawStdEncoding.DecodeString(sealed)
	if err != nil {
		return "", err
	}
	size := b.aead.NonceSize()
	if len(data) < size {
		return "", errors.New("seal: value too short")
	}
	plaintext, err := b.aead.Open(nil, data[:size], data[size:], nil)
	return string(plaintext), err
}
