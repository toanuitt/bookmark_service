package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	urlExpTime = 24 * time.Hour
)

// UrlStorage defines the interface for URL storage repository operations.
//
//go:generate mockery --name=UrlStorage --filename urlstorage_repo.go
type UrlStorage interface {
	StoreUrl(ctx context.Context, code, url string, exp int64) error
	GetUrl(ctx context.Context, code string) (string, error)
	Exist(ctx context.Context, code string) (bool, error)
}

// urlStorage is the concrete implementation of the UrlStorage interface.
type urlStorage struct {
	c *redis.Client
}

// NewUrlStorage creates and returns a new UrlStorage instance.
func NewUrlStorage(c *redis.Client) UrlStorage {
	return &urlStorage{c: c}
}

// StoreUrl stores a URL with its shortened code in Redis with a 24-hour expiration.
func (r *urlStorage) StoreUrl(ctx context.Context, code, url string, exp int64) error {
	return r.c.Set(ctx, code, url, urlExpTime).Err()
}

// GetUrl retrieves the original URL associated with the given shortened code.
func (r *urlStorage) GetUrl(ctx context.Context, code string) (string, error) {
	return r.c.Get(ctx, code).Result()
}

// Exist checks whether a shortened URL code exists in Redis.
func (r *urlStorage) Exist(ctx context.Context, code string) (bool, error) {
	result, err := r.c.Exists(ctx, code).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}
