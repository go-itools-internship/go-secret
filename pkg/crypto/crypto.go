/*
	Package crypto provides functions to encode and decode data
	Using aes crypto
*/
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"go.uber.org/zap"
)

type cryptographer struct {
	key        []byte
	randomFlag bool
	logger     *zap.SugaredLogger
}

func NewCryptographer(key []byte, logger *zap.SugaredLogger) *cryptographer {
	h := sha256.New()
	key32 := make([]byte, 32)
	copy(key32, h.Sum(key))
	return &cryptographer{
		key:        key32,
		randomFlag: true,
		logger:     logger,
	}
}

func (c *cryptographer) Encode(value []byte) ([]byte, error) {
	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("cryptographer, encode method: invalid key: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cryptographer, encode method: invalid size: %w", err)
	}
	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if c.randomFlag {
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, fmt.Errorf("cryptographer, encode method: unexpected data: %w", err)
		}
	}
	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case,
	// we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, value, nil)
	return ciphertext, nil
}

func (c *cryptographer) Decode(encodedValue []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("cryptographer, decode method: invalid key: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cryptographer, decode method: invalid size: %w", err)
	}
	// Get the nonce size
	nonceSize := aesGCM.NonceSize()
	// Extract the nonce from the encrypted data
	nonce, ciphertext := encodedValue[:nonceSize], encodedValue[nonceSize:]
	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("cryptographer, decode method: decryption error: %w", err)
	}
	return plaintext, nil
}
