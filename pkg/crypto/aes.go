// Package crypto provides encryption and decryption functions.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// AESCipher is an interface for AES encryption and decryption.
type AESCipher interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// NewAESCipher creates a new AESCipher.
func NewAESCipher(secretKey string) AESCipher {
	return &aesCipher{
		secretKey: secretKey,
	}
}

type aesCipher struct {
	secretKey string
}

// Encrypt encrypts a plaintext string.
func (a *aesCipher) Encrypt(plaintext string) (string, error) {
	cipherBlock, err := aes.NewCipher([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext), nil
}

// Decrypt decrypts a ciphertext string.
func (a *aesCipher) Decrypt(ciphertext string) (string, error) {
	cipherBlock, err := aes.NewCipher([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), err
}
