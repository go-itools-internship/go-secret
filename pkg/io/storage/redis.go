package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type redisVault struct {
	redisClient *redis.Client
	ctx         context.Context
}

func New(ctx context.Context, redisClient http.Client) (*redisVault, error) {
	rdb := redis.NewClient(&redis.Options{Addr: "", Password: "", DB: 0})
	rv := &redisVault{
		rdb, ctx,
	}
	return rv, nil
}

func (r *redisVault) SaveData(key, encodedValue []byte) error {
	err := r.redisClient.Set(r.ctx, string(key), encodedValue, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (r *redisVault) ReadData(key []byte) ([]byte, error) {
	val, err := r.redisClient.Get(r.ctx, string(key)).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}
