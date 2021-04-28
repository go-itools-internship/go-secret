package cmd

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoot_Set(t *testing.T) {
	t.Run("expect two keys", func(t *testing.T) {
		key := "key value"
		path := "testFile.txt"
		r := NewRoot()
		r.rootCmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipherKey", "ck", "--path", path})
		err := r.RootExecute(context.Background())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.Remove(path))
		}()

		testFile, err := os.Open(path)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, testFile.Close())
		}()

		fileData := make(map[string]string)
		require.NoError(t, json.NewDecoder(testFile).Decode(&fileData))

		var got string
		require.Len(t, fileData, 1)
		for key, _ := range fileData {
			got = key
			break // we iterate one time to get first key
		}
		require.EqualValues(t, key, got)
	})

	t.Run("expect two keys", func(t *testing.T) {
		firstKey := "first key"
		secondKey := "second key"
		path := "testFile.txt"
		r := NewRoot()
		r.rootCmd.SetArgs([]string{"set", "--key", firstKey, "--value", "test value", "--cipherKey", "ck", "--path", path})
		err := r.RootExecute(context.Background())
		require.NoError(t, err)

		r2 := NewRoot()
		r2.rootCmd.SetArgs([]string{"set", "--key", secondKey, "--value", "test value", "--cipherKey", "ck", "--path", path})
		err = r2.RootExecute(context.Background())
		require.NoError(t, err)

		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.Remove(path))
		}()

		testFile, err := os.Open(path)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, testFile.Close())
		}()

		fileData := make(map[string]string)
		require.NoError(t, json.NewDecoder(testFile).Decode(&fileData))

		require.Len(t, fileData, 2)
		require.Contains(t, fileData, firstKey)
		require.Contains(t, fileData, secondKey)
	})
}

func TestRoot_Get(t *testing.T) {
	key := "key value"
	value := "60OBdPOOkSOu6kn8ZuMuXtAPVrUEFkPREydDwY6+ip/LrAFaHSc="
	path := "testFile.txt"

	file, err := os.Create(path)
	require.NoError(t, err)
	fileTestData := make(map[string]string)
	fileTestData[key] = value
	require.NoError(t, json.NewEncoder(file).Encode(&fileTestData))

	defer func() {
		require.NoError(t, os.Remove(path))
	}()

	defer func() {
		require.NoError(t, file.Close())
	}()

	r := NewRoot()
	r.rootCmd.SetArgs([]string{"get", "--key", key, "--cipherKey", "ck", "--path", path})
	executeErr := r.RootExecute(context.Background())
	require.NoError(t, executeErr)

	testFile, err := os.Open(path)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, testFile.Close())
	}()

	fileData := make(map[string]string)
	require.NoError(t, json.NewDecoder(testFile).Decode(&fileData))
	var got string
	require.Len(t, fileData, 1)
	for _, value := range fileData {
		got = value // we iterate one time to get first value
		break
	}

	require.EqualValues(t, value, got)
}
