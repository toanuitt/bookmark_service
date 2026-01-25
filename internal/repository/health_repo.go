package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HealthCheckRepo defines the interface for health check repository operations.
//
//go:generate mockery --name=HealthCheckRepo  --filename health_repo.go
type HealthCheckRepo interface {
	Ping(ctx context.Context) error
}

// healthCheckStorage is the concrete implementation of the HealthCheckRepo interface.
type healthCheckStorage struct {
	redisClient *redis.Client
}

// NewHealthCheckStorage creates and returns a new HealthCheckRepo instance.
func NewHealthCheckStorage(redisClient *redis.Client) HealthCheckRepo {
	return &healthCheckStorage{redisClient: redisClient}
}

// Ping checks the health of the Redis connection.
func (r *healthCheckStorage) Ping(ctx context.Context) error {
	return r.redisClient.Ping(ctx).Err()
}
