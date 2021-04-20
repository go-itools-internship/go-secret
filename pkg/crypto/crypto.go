package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

type cryptographer struct {
	Key        []byte
	RandomFlag bool
}

func NewCryptographer(key []byte) *cryptographer {
	h := sha256.New()
	key32 := make([]byte, 32)
	copy(key32, h.Sum(key))
	return &cryptographer{
		Key:        key32,
		RandomFlag: true,
	}
}

func (c *cryptographer) Encode(value []byte) ([]byte, error) {
	// Create a new Cipher Block from the key
	// must be 16, 32, 64 bit key
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return nil, fmt.Errorf("i cryptograper, encode method: invalid key: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("i cryptograper, encode method: invalid size: %w", err)
	}
	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if c.RandomFlag {
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, fmt.Errorf("i cryptograper, encode method: unexpected data: %w", err)
		}
	}
	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case,
	// we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, value, nil)
	return ciphertext, nil
}

func (c *cryptographer) Decode(encodedValue []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.Key) // key
	if err != nil {
		return nil, fmt.Errorf("i cryptograper, decode method: invalid key: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("i cryptograper, decode method: invalid size: %w", err)
	}
	// Get the nonce size
	nonceSize := aesGCM.NonceSize()
	// Extract the nonce from the encrypted data
	nonce, ciphertext := encodedValue[:nonceSize], encodedValue[nonceSize:]
	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("i cryptograper, decode method: decryption error: %w", err)
	}
	return plaintext, nil
}
