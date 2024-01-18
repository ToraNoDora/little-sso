package redis_cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	ctx   context.Context
	cache *redis.Client
}

func NewRedisCache(cfg Config) *RedisCache {
	ctx := context.Background()
	ch, err := NewRedisClient(ctx, cfg)
	if err != nil {
		panic(err)
	}

	return &RedisCache{
		ctx:   ctx,
		cache: ch,
	}
}

func (rc *RedisCache) Set(key string, value interface{}) (bool, error) {
	result, err := rc.cache.SetNX(rc.ctx, key, value, 10*time.Second).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

func (rc *RedisCache) Get(key string) (interface{}, error) {
	val := rc.cache.Get(rc.ctx, key).Val()

	return val, nil
}
