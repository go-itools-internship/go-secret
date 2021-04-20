package crypto

import (
	"bytes"
	"fmt"
	"testing"
)

var value = []byte("All i need is love")
var key = []byte("I am the key")

func newCrypto(key []byte) *Cryptographer {
	key32 := make([]byte, 32)
	copy(key32, key)
	return &Cryptographer{
		Key:        key32,
		RandomFlag: false,
	}
}

func TestCryptographer_Encode(t *testing.T) {

	t.Run("encode", func(t *testing.T) {
		encode := newCrypto(key) // must 16, 32, 64 bit key
		got, _ := encode.Encode(value)
		var want []byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 123, 122, 21, 92, 176, 232, 107, 223, 7, 39, 146, 85, 75, 6, 225, 191, 12,
			95, 141, 227, 148, 9, 195, 36, 98, 150, 151, 193, 14, 63, 23, 152, 14, 35}
		fmt.Println([]byte(got))
		if !bytes.Equal(got, want) {
			t.Errorf(string(got), want)
		}
	})
}

func TestCryptographer_Decode(t *testing.T) {
	decode := newCrypto(key) // must 16, 32, 64 bit key
	got, _ := decode.Decode([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 123, 122, 21, 92, 176, 232, 107, 223, 7, 39, 146, 85, 75, 6,225, 191, 12,
		95, 141, 227, 148, 9, 195, 36, 98, 150, 151, 193, 14, 63, 23, 152, 14, 35})
	want := value

	if !bytes.Equal(got, want) {
		t.Errorf(string(got), want)
	}
}

