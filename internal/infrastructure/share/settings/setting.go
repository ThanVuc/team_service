package settings

import (
	"errors"
	"fmt"

	"github.com/thanvuc/go-core-lib/config"
)

func LoadConfiguration(path string) (*Configuration, error) {
	var cfg Configuration

	if err := config.LoadConfig(&cfg, path); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func validateConfig(cfg *Configuration) error {
	if cfg.Server.Host == "" {
		return errors.New("server.host is required")
	}

	if cfg.Redis.Host == "" {
		return errors.New("redis.host is required")
	}

	if cfg.RabbitMQ.Host == "" {
		return errors.New("rabbitmq.host is required")
	}

	return nil
}
