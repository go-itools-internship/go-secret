package storage

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

const testFilename = "testfile.json"

func TestFileVault(t *testing.T) {
	sugar := createSugarLogger()
	fileVault, err := NewFileVault(testFilename, sugar)
	require.NoError(t, err)

	defer func() {
		err = os.Remove(testFilename)
		if err != nil {
			sugar.Debug("file not removed")
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

func createSugarLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugar := logger.Sugar()
	return sugar
}
