package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var tests = []struct {
	name  string
	key   []byte
	value []byte
	want  string
}{
	{"encode/decode value 1", []byte("I am the key"), []byte("All i need is love"),
		"0000000000000000000000007b7a155cb0e86bdf072792554b06e1bf0c5f8de39409c324629697c10e3f17980e23"},
	{"encode/decode value 2", []byte("I am another key"), []byte("All i need is love love love"),
		"0000000000000000000000007658edf69909742db1371e66ab9dcfcec2fb449498c8326dd9c87674957a11b1d106d32ec5bde414664bac2a"},
	{"encode/decode with key match more than 32", []byte("werwewtwtwrtrtert55tttttttttttttggggggggggggrt56456hfghfhj$34g"), []byte("All i need is love"),
		"000000000000000000000000eff457b839858b568bcacf892b213d4591be73186c1ea6409f823ee4a4df806b95fc"},
	{"empty key and value", []byte(""), []byte(""),
		"000000000000000000000000530f8afbc74536b9a963b4f1c4cb738b"},
}

func TestCryptographer_Encode(t *testing.T) {
	for i, tt := range tests {
		t.Logf("\tTest: %d\tfor key %q and value %q", i+1, tt.key, tt.value)
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			encode := NewCryptographer(tt.key, &loopReader{})
			got, err := encode.Encode(tt.value)
			tgot := hex.EncodeToString(got)
			if err != nil {
				return
			}
			if tgot != tt.want {
				t.Errorf(string(got), tt.want)
			}
		})
	}
}

func TestCryptographer_Decode(t *testing.T) {
	for i, tt := range tests {
		t.Logf("\tTest: %d\tfor key %q and value %q", i+1, tt.key, tt.value)
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			decode := NewCryptographer(tt.key, &loopReader{})
			got, err := decode.Decode([]byte(tt.want))
			if err != nil {
				return
			}
			if !bytes.Equal(got, tt.value) {
				t.Errorf(string(got), tt.want)
			}
		})
	}
}
