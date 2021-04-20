package crypto

import (
	"bytes"
	"testing"
)

var value = []byte("All i need is love")
var key = []byte("I am the key")


func TestCryptographer_Encode(t *testing.T) {
		encode := NewCryptographer(key)// must 16, 32, 64 bit key
		encode.RandomFlag = false
		got, _ := encode.Encode(value)
		var want = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 123, 122, 21, 92, 176, 232, 107, 223, 7, 39, 146, 85, 75, 6, 225, 191, 12,
			95, 141, 227, 148, 9, 195, 36, 98, 150, 151, 193, 14, 63, 23, 152, 14, 35}
		if !bytes.Equal(got, want) {
			t.Errorf(string(got), want)
		}
}

func TestCryptographer_Decode(t *testing.T) {
	decode := NewCryptographer(key)// must 16, 32, 64 bit key
	decode.RandomFlag = false
	got, _ := decode.Decode([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 123, 122, 21, 92, 176, 232, 107, 223, 7, 39, 146, 85, 75, 6,225, 191, 12,
		95, 141, 227, 148, 9, 195, 36, 98, 150, 151, 193, 14, 63, 23, 152, 14, 35})
	want := value

	if !bytes.Equal(got, want) {
		t.Errorf(string(got), want)
	}
}
//
//func TestCryptographer_EncodeWithKeyMoreThan32(t *testing.T) {
//	encode := newCrypto([]byte("fkkfkkfkfkkfkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk")) // must 16, 32, 64 bit key
//	got, _ := encode.Encode(value)
//	want := make([]byte, 32)
//	if cap(encode.Key) != cap(want) {
//		t.Errorf(string(got), want)
//	}
//}
//
//func TestCryptographer_EncodeWithEmptyKey(t *testing.T) {
//	encode := newCrypto([]byte(""))
//	got, _ := encode.Encode(value)
//	want := make([]byte, 32)
//	if cap(encode.Key) != cap(want) {
//		t.Errorf(string(got), want)
//	}
//}

