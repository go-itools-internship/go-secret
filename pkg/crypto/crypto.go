/*
	Package crypto provides functions to encode and decode data
	Using aes crypto.
	Package crypto contain custom implementation of io Reader.
	Using for read data many times without changes.
*/
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"
)

type cryptographer struct {
	key         []byte
	nonceReader io.Reader
}

func NewCryptographer(key []byte, nonceReader io.Reader) *cryptographer {
	h := sha256.New()
	h.Write(key)
	key32 := make([]byte, 32)
	copy(key32, h.Sum(nil))
	return &cryptographer{
		key:         key32,
		nonceReader: nonceReader,
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

	if _, err = io.ReadFull(c.nonceReader, nonce); err != nil {
		return nil, fmt.Errorf("cryptographer, encode method: unexpected data: %w", err)
	}
	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case,
	// we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, value, nil)

	return ciphertext, nil
}

func (c *cryptographer) Decode(encodedValue []byte) ([]byte, error) {
	if encodedValue == nil {
		return nil, nil
	}
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
