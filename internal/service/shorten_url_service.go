package service

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/toanuitt/bookmark_service/internal/repository"
	"github.com/toanuitt/bookmark_service/pkg/stringutils"
)

const (
	maxRetryAttempts = 5
	LengthURLcode    = 7
)

var (
	ErrMaxRetriesExceeded = errors.New("Exceed for generate unique code URL")
	ErrURLNotFound        = errors.New("url not found")
)

// ShortenURLservice defines the interface for URL shortening operations.
//
//go:generate mockery --name ShortenURLservice --filename shorten_url_service.go
type ShortenURLservice interface {
	ShortlengthURL(ctx context.Context, originURL string, expireAt int64) (string, error)
	GetURL(ctx context.Context, url string) (string, error)
}

// shortenURL is the concrete implementation of the ShortenURLservice interface.
type shortenURL struct {
	repo    repository.UrlStorage
	codegen stringutils.CodeGenerator
}

// NewShortenURL creates and returns a new ShortenURLservice instance.
func NewShortenURL(repo repository.UrlStorage, codegen stringutils.CodeGenerator) ShortenURLservice {
	return &shortenURL{
		repo:    repo,
		codegen: codegen,
	}
}

// ShortlengthURL generates a unique shortened code for the given URL.
// It retries up to maxRetryAttempts times to ensure code uniqueness.
func (s *shortenURL) ShortlengthURL(ctx context.Context, originURL string, expireAt int64) (string, error) {
	var urlcode string
	for i := 1; i <= maxRetryAttempts; i++ {
		code, err := s.codegen.GenerateCode(LengthURLcode)
		if err != nil {
			return "", err
		}
		urlcode = code
		exists, err := s.repo.Exist(ctx, urlcode)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
		if i == maxRetryAttempts {
			return "", ErrMaxRetriesExceeded
		}
	}
	err := s.repo.StoreUrl(ctx, urlcode, originURL, expireAt)
	if err != nil {
		return "", err
	}
	return urlcode, nil
}

// GetURL retrieves the original URL for the given shortened code.
func (s *shortenURL) GetURL(ctx context.Context, url string) (string, error) {
	url, err := s.repo.GetUrl(ctx, url)
	if err != nil {
		if err == redis.Nil {
			return "", ErrURLNotFound
		}
		return "", err
	}
	return url, nil
}
