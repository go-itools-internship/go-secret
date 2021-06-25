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
		"000000000000000000000000a285210979aab1707d6215a0eba48236698b06fb4a20005ed0e5e24b538ca5e65107"},
	{"encode/decode value 2", []byte("I am another key"), []byte("All i need is love love love"),
		"000000000000000000000000a10c2db816d251e8242981a044d0452bd8abcad891a78a75ab0af64444318cbe4c28fe88f4acab3c4e347827"},
	{"encode/decode with key match more than 32", []byte("werwewtwtwrtrtert55tttttttttttttggggggggggggrt56456hfghfhj$34g"), []byte("All i need is love"),
		"000000000000000000000000918d52bfe6c8cd7898f5c2b7bd62b71ac34cdd8b177858493ac6184e23fe3188ea14"},
	{"empty key and value", []byte(""), []byte(""),
		"000000000000000000000000d51ed6081edb98739080fbe09ec476fb"},
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
