package storage

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const testFilename = "testfile.json"

func TestFileVault(t *testing.T) {
	fileVault, err := NewFileVault(testFilename)
	require.NoError(t, err)

	defer func() {
		err = os.Remove(testFilename)
		if err != nil {
			t.Logf("file not removed")
		}
	}()

	t.Run("SaveData", func(t *testing.T) {
		want := []byte("Hey")

		require.NoError(t, fileVault.SaveData([]byte("f1"), want))

		f, err := os.OpenFile(fileVault.path, os.O_RDONLY, 0600)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, f.Close())
		}()

		testStorage := make(map[string][]byte)
		require.NoError(t, json.NewDecoder(f).Decode(&testStorage))
		got := testStorage["f1"]

		require.EqualValues(t, want, got)
	})

	t.Run("ReadData", func(t *testing.T) {
		want := []byte("World")

		f, err := os.OpenFile(fileVault.path, os.O_WRONLY, 0600)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, f.Close())
		}()

		testStorage := make(map[string][]byte)
		testStorage["f2"] = want

		require.NoError(t, json.NewEncoder(f).Encode(testStorage))

		got, err := fileVault.ReadData([]byte("f2"))
		require.NoError(t, err)
		require.EqualValues(t, want, got)
	})
}
