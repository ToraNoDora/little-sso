package cache

import (
	rc "github.com/ToraNoDora/little-sso/sso/internal/src/store/cache/redis_cache"
)

type CacheStorage struct {
	Cache Cache
}

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) (bool, error)
}

func NewCache(cfg rc.Config) *CacheStorage {
	return &CacheStorage{
		Cache: rc.NewRedisCache(cfg),
	}
}
