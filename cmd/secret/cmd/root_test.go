package cmd

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRoot_Set(t *testing.T) {
	t.Run("expect one keys", func(t *testing.T) {
		key := "key value"
		path := "testFile.txt"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--path", path})
		err := r.Execute(ctx)
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
		for key := range fileData {
			got = key
			break // we iterate one time to get first key
		}
		require.EqualValues(t, key, got)
	})

	t.Run("expect two keys", func(t *testing.T) {
		firstKey := "first key"
		secondKey := "second key"
		path := "testFile.txt"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"set", "--key", firstKey, "--value", "test value", "--cipher-key", "ck", "--path", path})
		err := r.Execute(ctx)
		require.NoError(t, err)

		r2 := New()
		r2.cmd.SetArgs([]string{"set", "--key", secondKey, "--value", "test value", "--cipher-key", "ck", "--path", path})
		err = r2.Execute(ctx)
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(path))
	}()

	defer func() {
		require.NoError(t, file.Close())
	}()

	fileTestData := make(map[string]string)
	fileTestData[key] = value
	require.NoError(t, json.NewEncoder(file).Encode(&fileTestData))

	r := New()
	r.cmd.SetArgs([]string{"get", "--key", key, "--cipher-key", "ck", "--path", path})
	executeErr := r.Execute(ctx)
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
		got = value
		break // we iterate one time to get first value
	}
	require.EqualValues(t, value, got)
}

func TestRoot_Server(t *testing.T) {
	t.Run("expect success", func(t *testing.T) {

		key := "key value"
		path := "testFile.txt"
		port := "8888"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"server", "--cipher-key", key, "--path", path, "--port", port})
		var wg sync.WaitGroup
		wg.Add(1) // в группе две горутины
		work := func() {
			defer wg.Done()
			r := New()
			r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--path", path})
			err := r.Execute(ctx)
			require.NoError(t, err)

		}
		go work()
		wg.Wait()
		
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
		for key := range fileData {
			got = key
			break // we iterate one time to get first key
		}
		require.EqualValues(t, key, got)
	})

}
