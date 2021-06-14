package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRoot_Get(t *testing.T) {
	value := "60OBdPOOkSOu6kn8ZuMuXtAPVrUEFkPREydDwY6+ip/LrAFaHSc="
	t.Run("success", func(t *testing.T) {
		file, err := os.Create(path)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
	})
	t.Run("success after get redis command", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r := New()

		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--redis-url", redisURL})
		executeErr := r.Execute(ctx)
		require.NoError(t, executeErr)

		var b bytes.Buffer
		r.cmd.SetOut(&b)

		r.cmd.SetArgs([]string{"get", "--key", key, "--cipher-key", "ck", "--redis-url", redisURL})
		err := r.Execute(ctx)
		require.NoError(t, err)
		out := b.String()
		require.NoError(t, err)
		require.EqualValues(t, "test value\n", out)
	})
	t.Run("success after get postgres command", func(t *testing.T) {
		postgresURL := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r := New()

		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--postgres-url", postgresURL, "--migration", "file://../../../"})
		executeErr := r.Execute(ctx)
		require.NoError(t, executeErr)

		var b bytes.Buffer
		r.cmd.SetOut(&b)

		r.cmd.SetArgs([]string{"get", "--key", key, "--cipher-key", "ck", "--postgres-url", postgresURL})
		err := r.Execute(ctx)
		require.NoError(t, err)
		out := b.String()
		require.NoError(t, err)
		require.EqualValues(t, "test value\n", out)
	})
}
