package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/require"
)

func TestRedisVault_SaveData(t *testing.T) {
	key := "key"
	encodedValue := "value"
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		s := New(rdb, ctx)
		err := s.SaveData([]byte(key), []byte(encodedValue))
		require.NoError(t, err)
	})
	t.Run("error if key equals nil", func(t *testing.T) {
		key := ""
		s := New(rdb, ctx)
		err := s.SaveData([]byte(key), []byte(encodedValue))
		require.Error(t, err)
		require.EqualValues(t, "storage: key can't be nil ", err.Error())
	})
	t.Run("wrong key error if key has been deleted", func(t *testing.T) {
		nilEncodedValue := ""
		s := New(rdb, ctx)
		err := s.SaveData([]byte(key), []byte(key))
		require.NoError(t, err)
		err = s.SaveData([]byte(key), []byte(nilEncodedValue))
		require.NoError(t, err)

		val, err := s.ReadData([]byte(key))
		require.NoError(t, err)
		require.NotEqualValues(t, "", val)
	})
}

func TestRedisVault_ReadData(t *testing.T) {
	key := "key"
	encodedValue := "value"
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		s := New(rdb, ctx)

		err := s.SaveData([]byte(key), []byte(encodedValue))
		require.NoError(t, err)

		val, err := s.ReadData([]byte(key))
		require.NoError(t, err)
		fmt.Println(string(val))
		require.EqualValues(t, "value", val)
	})
	t.Run("error if wrong key", func(t *testing.T) {
		wrongKey := "wrongKey"
		s := New(rdb, ctx)
		val, err := s.ReadData([]byte(wrongKey))
		require.NoError(t, err)
		require.NotEqualValues(t, "value", val)
	})
}
