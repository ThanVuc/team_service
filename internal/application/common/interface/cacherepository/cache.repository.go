package icacherepository

import "context"

type CacheRepository interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl int) error
}
