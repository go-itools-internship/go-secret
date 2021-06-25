package storage

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type redisVault struct {
	client *redis.Client
}

// NewRedisVault create new redis client
func NewRedisVault(rdb *redis.Client) *redisVault {
	rv := &redisVault{
		client: rdb,
	}
	return rv
}

// SaveData put data in redis storage by key and encoded value
// 	key to set in redis storage (can't be nil)
// 	encoded value to storage
func (r *redisVault) SaveData(key, encodedValue []byte) error {
	ctx := context.Background()
	if bytes.Equal(key, []byte("")) {
		return errors.New("storage: key can't be nil")
	}
	if bytes.Equal(encodedValue, []byte("")) {
		fmt.Println("Key was deleted")
		_, err := r.client.Del(ctx, hex.EncodeToString(key)).Result()
		if err != nil {
			return fmt.Errorf("storage: %w", err)
		}
		return nil
	}
	fmt.Println(hex.EncodeToString(key))
	err := r.client.Set(ctx, hex.EncodeToString(key), encodedValue, 0).Err()
	if err != nil {
		return fmt.Errorf("storage: redis client can't set data %w", err)
	}
	return nil
}

// ReadData get data from redis storage by key
// 	key to get value for pair key-value from redis storage (can't be nil)
func (r *redisVault) ReadData(key []byte) ([]byte, error) {
	ctx := context.Background()
	if bytes.Equal(key, []byte("")) {
		return nil, errors.New("storage: key can't be nil")
	}
	val, err := r.client.Get(ctx, hex.EncodeToString(key)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("storage: redis client can't get data %w", err)
	}
	return []byte(val), nil
}
