package storage

import (
	"fmt"
	"testing"
)

func TestFileVault_SaveData(t *testing.T) {

	fv, err := NewFileVault("file")
	if err != nil {
		fmt.Println(err)
		t.Error("newFileVault = nil")
	}
	got := []byte("Hello")

	err = fv.SaveData([]byte("f1"), got)
	if err != nil {
		t.Error("file not found or something wrong")
	}

	wont := fv.storage["f1"]
	if string(wont) != string(got) {
		t.Error("file not found or something wrong")
	}
}

func TestFileVault_ReadData(t *testing.T) {

	fv, err := NewFileVault("file")
	if err != nil {
		fmt.Println(err)
		t.Error("newFileVault = nil")
	}
	wont := []byte("World")
	fv.storage["f2"] = wont

	got, err := fv.ReadData([]byte("f2"))
	if err != nil || string(got) != string(wont) {
		t.Error("file not found or something wrong")
	}
}
