package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/go-itools-internship/go-secret/pkg/io/storage"

	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/require"
)

func TestRoot_Set(t *testing.T) {
	t.Run("expect one keys", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
		require.NotEmpty(t, got)
	})
	t.Run("expect set data only redis storage", func(t *testing.T) {
		key := "12345"
		cipherKey := "d93beca6efd0421b314c08102f48f13038410738dbd2050c04fc89265a10024ba4"
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		r := New()

		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--redis-url", redisURL})
		err := r.Execute(ctx)
		require.NoError(t, err)

		rdb := redis.NewClient(&redis.Options{Addr: redisURL, Password: "", DB: 0})

		val, err := rdb.Get(ctx, cipherKey).Result()
		require.NoError(t, err)
		require.NotEmpty(t, val)
	})
	t.Run("expect success set data postgres storage if get key error", func(t *testing.T) {
		key := "12345"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		defer func() {
			err := migrateDown(t)
			if err != nil {
				fmt.Println("can't migrate down ", err)
			}
		}()
		r := New()

		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--postgres-url", postgresURL, "--migration", migration})
		err := r.Execute(ctx)
		require.NoError(t, err)

		db, err := sqlx.ConnectContext(ctx, "postgres", postgresURL)
		defer disconnectPDB(db, r.logger.Named("test"))
		require.NoError(t, err)

		d := storage.NewPostgreVault(db)
		_, err = d.ReadData([]byte("12345"))
		require.Error(t, err)
		require.EqualValues(t, err.Error(), "postgres: key not found")
	})

	t.Run("expect two keys", func(t *testing.T) {
		firstKey := "first key"
		secondKey := "second key"
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
		require.NotEmpty(t, fileData)
		require.Len(t, fileData, 2)
	})
}
