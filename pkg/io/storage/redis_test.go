package storage

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/require"
)

func TestRedisVault_SaveData(t *testing.T) {
	key := "key"
	encodedValue := "value"
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer disconnectRDB(rdb, t)
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s := NewRedisVault(rdb)
		err := s.SaveData([]byte(key), []byte(encodedValue))
		require.NoError(t, err)

		val, err := s.client.Get(ctx, key).Result()
		require.NoError(t, err)
		require.EqualValues(t, encodedValue, val)
	})
	t.Run("error if key equals nil", func(t *testing.T) {
		key := ""
		s := NewRedisVault(rdb)
		err := s.SaveData([]byte(key), []byte(encodedValue))
		require.Error(t, err)
		require.EqualValues(t, "storage: key can't be nil ", err.Error())
	})
	t.Run("get nil value if key has been deleted", func(t *testing.T) {
		nilEncodedValue := ""
		s := NewRedisVault(rdb)
		err := s.SaveData([]byte(key), []byte(key))
		require.NoError(t, err)
		err = s.SaveData([]byte(key), []byte(nilEncodedValue))
		require.NoError(t, err)

		val, err := s.ReadData([]byte(key))
		require.NoError(t, err)
		require.EqualValues(t, []byte(nil), val)
	})
}

func TestRedisVault_ReadData(t *testing.T) {
	key := "key"
	encodedValue := "value"
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer disconnectRDB(rdb, t)
	t.Run("success", func(t *testing.T) {
		s := NewRedisVault(rdb)
		err := s.SaveData([]byte(key), []byte(encodedValue))
		require.NoError(t, err)
		val, err := s.ReadData([]byte(key))
		require.NoError(t, err)
		require.EqualValues(t, "value", val)
	})
	t.Run("get nil value if wrong key", func(t *testing.T) {
		wrongKey := "wrongKey"
		s := NewRedisVault(rdb)
		val, err := s.ReadData([]byte(wrongKey))
		require.NoError(t, err)
		require.EqualValues(t, []byte(nil), val)
	})
}

func disconnectRDB(rdb *redis.Client, t *testing.T) {
	err := rdb.Close()
	if err != nil {
		t.Log("can't disconnect redis db")
	}
}
