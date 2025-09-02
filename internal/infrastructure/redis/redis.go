package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"medika-backend/internal/infrastructure/config"
)

func New(cfg config.RedisConfig) (*redis.Client, error) {
	// Parse Redis URL
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Configure connection settings
	opt.MaxRetries = cfg.MaxRetries
	opt.MinIdleConns = cfg.MinIdleConns
	opt.PoolSize = cfg.PoolSize
	opt.ReadTimeout = cfg.ReadTimeout
	opt.WriteTimeout = cfg.WriteTimeout

	// Create Redis client
	client := redis.NewClient(opt)

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}
