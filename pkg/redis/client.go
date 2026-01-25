package redis

import (
	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and returns a new Redis client instance.
// It reads configuration from environment variables and initializes the client.
func NewRedisClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return client, nil

}
