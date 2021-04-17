package crypto

import (
	"crypto/aes"
	"encoding/hex"
	"github.com/go-itools-internship/go-secret/pkg/secret"
)

type Cryptographer struct {
	crypto *secret.Cryptographer
	Key    []byte
}

func  encodeHex(b []byte) []byte {
	return []byte(hex.EncodeToString(b))
}

func  decodeHex(b []byte) []byte {
	data, err := hex.DecodeString(string(b))
	CheckError(err)
	return data
}

func (c *Cryptographer) Encode(value []byte) ([]byte, error) {
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(c.Key)
	CheckError(err)
	// allocate space for ciphered data
	out := make([]byte, len(value))
	// encrypt
	block.Encrypt(out, value)
	return encodeHex(out), nil
}

func (c *Cryptographer) Decode(encodedValue []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.Key) // key?
	CheckError(err)

	out := make([]byte, len(encodedValue))
	block.Decrypt(out, decodeHex(encodedValue))

	return out, nil
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
