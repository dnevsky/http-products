package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisCache(cfg *redis.Options) (*redis.Client, error) {
	cache := redis.NewClient(cfg)

	res, err := cache.Ping(context.Background()).Result()

	fmt.Println(res)

	if err != nil {
		return nil, err
	}

	return cache, nil
}
