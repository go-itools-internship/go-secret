package crypto

import (
	"bytes"
	"testing"
)

var tests = []struct {
	name  string
	key   []byte
	value []byte
	want  []byte
}{
	{"encode/decode value 1", []byte("I am the key"), []byte("All i need is love"),
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 103, 219, 30, 14, 204, 218, 207, 85, 148, 66,
			160, 22, 65, 82, 133, 234, 239, 51, 104, 63, 206, 168, 142, 30, 10, 255, 243, 84, 85, 36, 201, 78, 183, 51}},
	{"encode/decode value 2", []byte("I am another key"), []byte("All i need is love love love"),
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 170, 62, 141, 236, 82, 152, 165, 131, 27, 42, 227, 232, 125, 148,
			125, 23, 77, 242, 99, 196, 24, 2, 221, 82, 234, 5, 16, 73, 28, 153, 134, 58,
			197, 134, 198, 43, 27, 28, 162, 145, 82, 157, 6, 121}},
	{"encode/decode with key match more than 32", []byte("werwewtwtwrtrtert55tttttttttttttggggggggggggrt56456hfghfhj$34g"), []byte("All i need is love"),
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 239, 244, 87, 184, 57, 133, 139, 86, 139, 202, 207, 137, 43, 33, 61, 69,
			145, 190, 115, 24, 108, 30, 166, 64, 159, 130, 62, 228, 164, 223, 128, 107, 149, 252}},
	{"empty key and value", []byte(""), []byte(""),
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 213, 30, 214, 8, 30, 219, 152, 115, 144, 128, 251, 224, 158, 196, 118, 251}},
}

func TestCryptographer_Encode(t *testing.T) {
	for i, tt := range tests {
		t.Logf("\tTest: %d\tfor key %q and value %q", i+1, tt.key, tt.value)
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			encode := NewCryptographer(tt.key) // must 16, 32, 64 bit key
			encode.RandomFlag = false
			got, err := encode.Encode(tt.value)
			if err != nil {
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf(string(got), tt.want)
			}
		})
	}
}

func TestCryptographer_Decode(t *testing.T) {
	for i, tt := range tests {
		t.Logf("\tTest: %d\tfor key %q and value %q", i+1, tt.key, tt.value)
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			decode := NewCryptographer(tt.key) // must 16, 32, 64 bit key
			decode.RandomFlag = false
			got, err := decode.Decode(tt.want)
			if err != nil {
				return
			}
			if !bytes.Equal(got, tt.value) {
				t.Errorf(string(got), tt.want)
			}
		})
	}
}
