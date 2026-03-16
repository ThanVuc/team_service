package cacherepository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thanvuc/go-core-lib/cache"
)

type CacheRepository struct {
	cache  *cache.RedisCache
	client *redis.Client
}

func NewCacheRepository(cache *cache.RedisCache) *CacheRepository {
	return &CacheRepository{
		cache:  cache,
		client: cache.Client,
	}
}

func (r *CacheRepository) Get(ctx context.Context, key string, dest any) error {
	return r.client.Get(ctx, key).Scan(dest)
}

func (r *CacheRepository) Set(ctx context.Context, key string, value any, ttl int) error {
	timeDuration := time.Duration(ttl) * time.Second
	return r.client.Set(ctx, key, value, timeDuration).Err()
}
