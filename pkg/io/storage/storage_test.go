package storage

import (
	"fmt"
	"testing"
)

func TestFileVault_SaveData(t *testing.T) {

	fileVault, err := NewFileVault("file")
	if err != nil {
		fmt.Println(err)
		t.Error("fileVault = nil")
	}
	got := []byte("Hey")

	err = fileVault.SaveData([]byte("f1"), got)
	if err != nil {
		t.Error("file not found or something wrong")
	}

	defer func() {
		err := fileVault.Close()
		if err != nil {
			t.Error("file not closed")
		}
	}()

	wont := fileVault.storage["f1"]
	if string(wont) != string(got) || err != nil {
		t.Error("file not found or something wrong")
	}
}

func TestFileVault_ReadData(t *testing.T) {

	fileVault, err := NewFileVault("file")
	if err != nil {
		fmt.Println(err)
		t.Error("fileVault = nil")
	}
	wont := []byte("World")
	fileVault.storage["f2"] = wont

	got, err := fileVault.ReadData([]byte("f2"))
	if err != nil || string(got) != string(wont) {
		t.Error("file not found or something wrong")
	}

	defer func() {
		err := fileVault.Close()
		if err != nil {
			t.Error("file not closed")
		}
	}()
}
