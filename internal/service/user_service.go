// Package service provides business logic and operations for the bookmark service.
package service

import (
	"context"
	"time"

	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/internal/repository"
	"github.com/toanuitt/bookmark_service/pkg/utils"
)

//go:generate mockery --name Userservice --filename user_service.go

// Userservice defines business logic related to users.
type Userservice interface {
	Register(ctx context.Context, username, password, displayName, email string) (*model.User, error)
}

type user struct {
	repo repository.UserRepo
}

// NewUser creates a new Userservice instance.
func NewUser(repo repository.UserRepo) Userservice {
	return &user{repo: repo}
}

// Register registers a new user.
// It hashes the password, sets timestamps, and persists the user via the repository.
func (s *user) Register(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	hashedPwd := utils.HashPassword(password)
	now := time.Now()

	newUser := &model.User{
		Username:    username,
		Password:    hashedPwd,
		DisplayName: displayName,
		Email:       email,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	res, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return res, nil
}
