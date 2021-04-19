package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type Cryptographer struct {
	Key    []byte
}

func NewCryptographer(key []byte) *Cryptographer{
	return &Cryptographer{
		Key: key,
	}
}

func encodeHex(b []byte) []byte {
	return []byte(hex.EncodeToString(b))
}

func decodeHex(b []byte) []byte {
	data, err := hex.DecodeString(string(b))
	CheckError(err)
	return data
}

func encodeGCM(block cipher.Block, value []byte) []byte {
	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	CheckError(err)
	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, value, nil)
	return ciphertext
}

func decodeGCM(block cipher.Block, value []byte) []byte {
	aesGCM, err := cipher.NewGCM(block)
	CheckError(err)
	//Get the nonce size
	nonceSize := aesGCM.NonceSize()
	//Extract the nonce from the encrypted data
	nonce, ciphertext := value[:nonceSize], value[nonceSize:]
	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	CheckError(err)
	return plaintext
}

func (c *Cryptographer) Encode(value []byte) ([]byte, error) {
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(c.Key)
	CheckError(err)
	// allocate space for ciphered data
	plaintext := make([]byte, len(value))
	// encrypt
	block.Encrypt(plaintext, value)
	return encodeHex(plaintext), nil
}

func (c *Cryptographer) Decode(encodedValue []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.Key) // key?
	CheckError(err)

	plaintext := make([]byte, len(encodedValue))
	block.Decrypt(plaintext, decodeHex(encodedValue))

	return plaintext, nil
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
