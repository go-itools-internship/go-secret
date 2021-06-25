package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRoot_Get(t *testing.T) {
	value := "2Tvspu/QQhsxTAgQah+xcC3VhifWUjlZHKJYgYIabyiTK4Bx6Zo="
	t.Run("success", func(t *testing.T) {
		file, err := os.Create(path)
		require.NoError(t, err)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--path", path})
		err = r.Execute(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.Remove(path))
		}()

		defer func() {
			require.NoError(t, file.Close())
		}()

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
	t.Run("error after get file command with wrong ck", func(t *testing.T) {
		file, err := os.Create(path)
		require.NoError(t, err)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--path", path})
		err = r.Execute(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.Remove(path))
		}()

		defer func() {
			require.NoError(t, file.Close())
		}()

		var b bytes.Buffer
		r.cmd.SetOut(&b)

		r.cmd.SetArgs([]string{"get", "--key", key, "--cipher-key", "wrong-ck", "--path", path})
		err = r.Execute(ctx)
		require.Error(t, err)
		out := b.String()
		require.EqualValues(t, "can't get data by key: provider, GetData method: read data error: filevault: cannot read data: not found", err.Error())
		require.Empty(t, out)
	})
	t.Run("success after get redis command", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		defer func() {
			fmt.Println("postgres test: try migrate down")
			err := migrateDown(t)
			if err != nil {
				fmt.Println("can't migrate down", err)
			}
		}()

		r := New()

		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--postgres-url", postgresURL, "--migration", migration})
		executeErr := r.Execute(ctx)
		require.NoError(t, executeErr)

		var b bytes.Buffer
		r.cmd.SetOut(&b)

		r.cmd.SetArgs([]string{"get", "--key", key, "--cipher-key", "ck", "--postgres-url", postgresURL, "--migration", migration})
		err := r.Execute(ctx)
		require.NoError(t, err)
		out := b.String()
		require.NoError(t, err)
		require.EqualValues(t, "test value\n", out)
	})
	t.Run("error after get redis command with wrong ck", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		r := New()

		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--redis-url", redisURL})
		executeErr := r.Execute(ctx)
		require.NoError(t, executeErr)

		var b bytes.Buffer
		r.cmd.SetOut(&b)

		r.cmd.SetArgs([]string{"get", "--key", key, "--cipher-key", "wrong-ck", "--redis-url", redisURL})
		err := r.Execute(ctx)
		require.NoError(t, err)
		out := b.String()
		require.EqualValues(t, "\n", out)
	})
	t.Run("error after get postgres command with wrong ck", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		defer func() {
			fmt.Println("postgres test: try migrate down")
			err := migrateDown(t)
			if err != nil {
				fmt.Println("can't migrate down", err)
			}
		}()

		r := New()

		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--postgres-url", postgresURL, "--migration", migration})
		executeErr := r.Execute(ctx)
		require.NoError(t, executeErr)

		var b bytes.Buffer
		r.cmd.SetOut(&b)

		r.cmd.SetArgs([]string{"get", "--key", key, "--cipher-key", "wrong-ck", "--postgres-url", postgresURL, "--migration", migration})
		err := r.Execute(ctx)
		require.Error(t, err)
		out := b.String()
		require.Empty(t, out)
		require.EqualValues(t, "can't get data by key: provider, GetData method: read data error: postgres: key not found", err.Error())
	})
}
