package infracache

import (
	"fmt"
	"team_service/internal/infrastructure/share/settings"

	"github.com/thanvuc/go-core-lib/cache"
	"github.com/thanvuc/go-core-lib/log"
)

func NewRedisClient(cfg settings.Redis, logger log.LoggerV2) (*cache.RedisCache, error) {

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	client := cache.NewRedisCache(cache.Config{
		Addr:     addr,
		DB:       cfg.DB,
		Password: cfg.Password,
		PoolSize: cfg.PoolSize,
		MinIdle:  cfg.MinIdle,
	})

	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	logger.Info(fmt.Sprintf("Redis client connected to %s", addr))

	return client, nil
}
