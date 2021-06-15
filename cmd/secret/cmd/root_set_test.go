package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-itools-internship/go-secret/pkg/io/storage"
	"github.com/jmoiron/sqlx"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

func TestRoot_Set(t *testing.T) {
	t.Run("expect one keys", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
	t.Run("expect set data only redis storage", func(t *testing.T) {
		key := "12345"
		path := ""
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--redis-url", redisURL})
		err := r.Execute(ctx)
		require.NoError(t, err)

		_, err = os.Open(path)
		require.Error(t, err)

		rdb := redis.NewClient(&redis.Options{Addr: redisURL, Password: "", DB: 0})

		val, err := rdb.Get(ctx, key).Result()
		require.NoError(t, err)
		require.NotEmpty(t, val)
	})
	t.Run("expect set data only postgres storage", func(t *testing.T) {
		key := "12345"
		path := ""
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		defer func() {
			err := migrateDown()
			if err != nil {
				fmt.Println("can't migrate down ", err)
			}
		}()
		r := New()
		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--postgres-url", postgresURL, "--migration", migration})
		err := r.Execute(ctx)
		require.NoError(t, err)

		_, err = os.Open(path)
		require.Error(t, err)

		db, err := sqlx.ConnectContext(ctx, "postgres", postgresURL)
		require.NoError(t, err)

		d := storage.NewPostgreVault(db)
		data, err := d.ReadData([]byte("12345"))
		require.NoError(t, err)
		require.NotEmpty(t, data)
	})

	t.Run("expect two keys", func(t *testing.T) {
		firstKey := "first key"
		secondKey := "second key"
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
