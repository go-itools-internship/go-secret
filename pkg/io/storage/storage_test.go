package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

		testStorage := make(map[string][]byte)
		err = json.NewDecoder(file).Decode(&testStorage)
		if err != nil {
			t.Error("unable to write data from file to map")
		}
		got := testStorage["f1"]

		assert.EqualValues(t, want, got)
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

		fileVault.storage["f2"] = want

		err = json.NewEncoder(file).Encode(fileVault.storage)
		if err != nil {
			t.Error("unable to write data from map to file")
		}

		got, err := fileVault.ReadData([]byte("f2"))
		if err != nil || string(got) != string(want) {
			t.Error("file not found or something wrong")
		}

		assert.EqualValues(t, want, got)
	})

	err = os.Remove("file")
	if err != nil {
		t.Error("file not removed")
	}
}
