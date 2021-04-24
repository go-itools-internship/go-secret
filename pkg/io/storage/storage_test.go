package storage

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestFileVault(t *testing.T) {

	fileVault, err := NewFileVault("file")
	if err != nil {
		t.Error("fileVault = nil")
	}

	t.Run("SaveData", func(t *testing.T) {
		want := []byte("Hey")

		err = fileVault.SaveData([]byte("f1"), want)
		if err != nil {
			t.Error("file not found or something wrong")
		}

		file, err := os.OpenFile(fileVault.path, os.O_RDONLY, 0600)
		if err != nil {
			t.Error("unable to open file")
		}
		defer func() {
			cerr := file.Close()
			if err == nil {
				err = cerr
			}
		}()
		err = json.NewDecoder(file).Decode(&fileVault.Storage)
		if err != nil {
			t.Error("unable to write data from file to map")
		}
		got := fileVault.Storage["f1"]

		if !bytes.Equal(got, want) {
			t.Errorf("valies is different")
		}
	})

	t.Run("ReadData", func(t *testing.T) {
		want := []byte("World")

		file, err := os.OpenFile(fileVault.path, os.O_WRONLY, 0600)
		if err != nil {
			t.Error("unable to open file")
		}
		defer func() {
			cerr := file.Close()
			if err == nil {
				err = cerr
			}
		}()

		fileVault.Storage["f2"] = want

		err = json.NewEncoder(file).Encode(fileVault.Storage)
		if err != nil {
			t.Error("unable to write data from map to file")
		}

		got, err := fileVault.ReadData([]byte("f2"))
		if err != nil || string(got) != string(want) {
			t.Error("file not found or something wrong")
		}

		if !bytes.Equal(got, want) {
			t.Errorf("valies is different")
		}
	})

	err = os.Remove("file")
	if err != nil {
		t.Error("file not removed")
	}
}
