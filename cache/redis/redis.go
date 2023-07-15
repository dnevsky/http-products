package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

func NewRedisCache(cfg *redis.Options) (*redis.Client, error) {
	cache := redis.NewClient(cfg)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()

	res, err := cache.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}

	if res != "PONG" {
		return nil, errors.New("failed connect to Redis")
	}

	return cache, nil
}
