package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	urlExpTime = 24 * time.Hour
)

//go:generate mockery --name=UrlStorage --filename urlstorage_repo.go
type UrlStorage interface {
	StoreUrl(ctx context.Context, code, url string, exp int) error
	GetUrl(ctx context.Context, code string) (string, error)
	Exist(ctx context.Context, code string) (bool, error)
}

type urlStorage struct {
	c *redis.Client
}

func NewUrlStorage(c *redis.Client) UrlStorage {
	return &urlStorage{c: c}
}

func (r *urlStorage) StoreUrl(ctx context.Context, code, url string, exp int) error {
	return r.c.Set(ctx, code, url, urlExpTime).Err()
}

func (r *urlStorage) GetUrl(ctx context.Context, code string) (string, error) {
	return r.c.Get(ctx, code).Result()
}

func (r *urlStorage) Exist(ctx context.Context, code string) (bool, error) {
	result, err := r.c.Exists(ctx, code).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}
