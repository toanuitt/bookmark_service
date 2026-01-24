package service

import (
	"context"
	"errors"

	"github.com/toanuitt/bookmark_service/internal/repository"
	"github.com/toanuitt/bookmark_service/pkg/stringutils"
)

const (
	maxRetryAttempts = 5
	LengthURLcode    = 7
)

//go:generate mockery --name ShortenURLservice --filename shorten_url_service.go
type ShortenURLservice interface {
	ShortlengthURL(ctx context.Context, originURL string, expireAt int) (string, error)
	GetURL(ctx context.Context, url string) (string, error)
}

type shortenURL struct {
	repo    repository.UrlStorage
	codegen stringutils.CodeGenerator
}

func NewShortenURL(repo repository.UrlStorage, codegen stringutils.CodeGenerator) ShortenURLservice {
	return &shortenURL{
		repo:    repo,
		codegen: codegen,
	}
}

func (s *shortenURL) ShortlengthURL(ctx context.Context, originURL string, expireAt int) (string, error) {
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
		// If code doesn't exist, we can use it
		if !exists {
			break
		}
		// If this is the last attempt and code still exists, fail
		if i == maxRetryAttempts {
			return "", errors.New("Exceed for generate unique code URL")
		}
	}
	err := s.repo.StoreUrl(ctx, urlcode, originURL, expireAt)
	if err != nil {
		return "", err
	}
	return urlcode, nil
}

func (s *shortenURL) GetURL(ctx context.Context, url string) (string, error) {
	url, err := s.repo.GetUrl(ctx, url)
	if err != nil {
		return "", err
	}
	return url, nil
}
