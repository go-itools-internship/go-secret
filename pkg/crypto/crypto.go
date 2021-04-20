package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type Cryptographer struct {
	Key        []byte
	RandomFlag bool
}

func NewCryptographer(key []byte) *Cryptographer {
	key32 := make([]byte, 32)
	copy(key32, key)
	return &Cryptographer{
		Key:        key32,
		RandomFlag: true,
	}
}

func (c *Cryptographer) Encode(value []byte) ([]byte, error) {
	//Create a new Cipher Block from the key
	//must be 16, 32, 64 bit key
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		fmt.Println("Invalid key", err)
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Invalid size", err)
		return nil, err
	}
	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if c.RandomFlag {
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			fmt.Println("Unexpected data", err)
			return nil, err
		}
	}
	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case,
	//we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, value, nil)
	return ciphertext, nil
}

func (c *Cryptographer) Decode(encodedValue []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.Key) // key
	if err != nil {
		fmt.Println("Invalid key", err)
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Invalid size", err)
		return nil, err
	}
	//Get the nonce size
	nonceSize := aesGCM.NonceSize()
	//Extract the nonce from the encrypted data
	nonce, ciphertext := encodedValue[:nonceSize], encodedValue[nonceSize:]
	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("Decryption error ", err)
		return nil, err
	}
	return plaintext, nil
}
