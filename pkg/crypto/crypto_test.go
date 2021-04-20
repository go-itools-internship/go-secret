package crypto

import (
	"bytes"
	"fmt"
	"testing"
)

var tests = []struct {
	name  string
	key   []byte
	value []byte
	want  []byte
}{
	{"encode/decode value 1", []byte("I am the key"), []byte("All i need is love"), []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 123, 122, 21, 92, 176, 232,
		107, 223, 7, 39, 146, 85, 75, 6, 225, 191, 12,
		95, 141, 227, 148, 9, 195, 36, 98, 150, 151, 193, 14, 63, 23, 152, 14, 35}},
	{"encode/decode value 2", []byte("I am another key"), []byte("All i need is love love love"), []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 118, 88, 237, 246, 153, 9, 116,
		45, 177, 55, 30, 102, 171, 157, 207, 206, 194, 251,
		68, 148, 152, 200, 50, 109, 217, 200, 118, 116, 149, 122, 17, 177, 209, 6, 211, 46, 197, 189, 228, 20, 102, 75, 172, 42}},
	{"encode/decode with key match more than 32", []byte("werwewtwtwrtrtert55tttttttttttttggggggggggggrt56456hfghfhj$34g"), []byte("All i need is love"),
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 239, 244, 87, 184, 57, 133, 139, 86, 139, 202, 207, 137, 43, 33, 61, 69,
			145, 190, 115, 24, 108, 30, 166, 64, 159, 130, 62, 228, 164, 223, 128, 107, 149, 252}},
	{"empty key and value", []byte(""), []byte(""),
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 83, 15, 138, 251, 199, 69, 54, 185, 169, 99, 180, 241, 196, 203, 115, 139}},
}

func TestCryptographer_Encode(t *testing.T) {
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("\tTest: %d\tfor key %q and value %q", i+1, tt.key, tt.value)
			{
				encode := NewCryptographer(tt.key) // must 16, 32, 64 bit key
				encode.RandomFlag = false
				got, _ := encode.Encode(tt.value)
				fmt.Println(got)
				if !bytes.Equal(got, tt.want) {
					t.Errorf(string(got), tt.want)
				}
			}
		})
	}
}

func TestCryptographer_Decode(t *testing.T) {
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("\tTest: %d\tfor key %q and value %q", i+1, tt.key, tt.value)
			{
				decode := NewCryptographer(tt.key) // must 16, 32, 64 bit key
				decode.RandomFlag = false
				got, _ := decode.Decode(tt.want)
				if !bytes.Equal(got, tt.value) {
					t.Errorf(string(got), tt.want)
				}
			}
		})
	}
}

