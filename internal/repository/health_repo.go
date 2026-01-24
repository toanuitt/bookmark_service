package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=HealthCheckRepo  --filename health_repo.go
type HealthCheckRepo interface {
	Ping(ctx context.Context) error
}

type healthCheckStorage struct {
	redisClient *redis.Client
}

func NewHealthCheckStorage(redisClient *redis.Client) HealthCheckRepo {
	return &healthCheckStorage{redisClient: redisClient}
}

func (r *healthCheckStorage) Ping(ctx context.Context) error {
	return r.redisClient.Ping(ctx).Err()
}
