package urlstorage

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
//go:generate mockery --name=UrlStorage --filename urlstorage.go
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
