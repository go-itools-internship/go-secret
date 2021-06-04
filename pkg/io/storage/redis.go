package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type redisVault struct {
	redisClient *redis.Client
	ctx         context.Context
}

// New create new redis client with context
func New(rdb *redis.Client, ctx context.Context) *redisVault {
	rv := &redisVault{
		rdb, ctx,
	}
	return rv
}

// SaveData put data in redis storage by key and encoded value
// 	key to set in redis storage (can't be nil)
// 	encoded value to storage
func (r *redisVault) SaveData(key, encodedValue []byte) error {
	if bytes.Equal(key, []byte("")) {
		return errors.New("storage: key can't be nil ")
	}
	if bytes.Equal(encodedValue, []byte("")) {
		fmt.Println("Key was deleted")
		res, err := r.redisClient.Del(r.ctx, string(key)).Result()
		fmt.Println(res)
		if err != nil {
			return fmt.Errorf("storage: %w", err)
		}
		return nil
	}
	err := r.redisClient.Set(r.ctx, string(key), encodedValue, 0).Err()
	if err != nil {
		return fmt.Errorf("storage: redis client can't set data %w", err)
	}
	return nil
}

// ReadData get data from redis storage by key
// 	key to get value for pair key-value from redis storage (can't be nil)
func (r *redisVault) ReadData(key []byte) ([]byte, error) {
	if bytes.Equal(key, []byte("")) {
		return nil, errors.New("storage: key can't be nil ")
	}
	val, err := r.redisClient.Get(r.ctx, string(key)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("storage: redis client can't get data %w", err)
	}
	return []byte(val), nil
}
