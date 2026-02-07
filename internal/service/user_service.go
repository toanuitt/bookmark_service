package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/internal/repository"
	"github.com/toanuitt/bookmark_service/pkg/jwtutils"
	"github.com/toanuitt/bookmark_service/pkg/utils"
)

//go:generate mockery --name Userservice --filename user_service.go

const (
	tokenLast = 24 * time.Hour
)

var (
	// ErrClientErr is returned when login credentials are invalid.
	ErrClientErr = errors.New("invalid username or password")

	// ErrClientNoUpdate is returned when UpdateUser is called with no fields to update.
	ErrClientNoUpdate = errors.New("no update")
)

// Userservice defines business logic related to users.
type Userservice interface {
	Register(ctx context.Context, username, password, displayName, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
	UpdateUser(ctx context.Context, userID, displayName, email string) error
}

type user struct {
	repo   repository.UserRepo
	jwtGen jwtutils.JWTGenerator
}

// NewUser creates a new Userservice instance.
func NewUser(repo repository.UserRepo, jwtGen jwtutils.JWTGenerator) Userservice {
	return &user{repo: repo, jwtGen: jwtGen}
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

func (s *user) Login(ctx context.Context, username, password string) (string, error) {
	// check if user exist
	chosenUser, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	// check if password is valid
	check := utils.VerifyPassword(password, chosenUser.Password)
	if !check {
		return "", ErrClientErr
	}

	// create token
	jwtContent := jwt.MapClaims{
		"sub": chosenUser.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tokenLast).Unix(),
	}
	token, err := s.jwtGen.GenerateToken(jwtContent)
	if err != nil {
		return "", err
	}

	// return token
	return token, nil
}

func (s *user) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	return s.repo.GetUserById(ctx, userId)
}

func (s *user) UpdateUser(ctx context.Context, userID, displayName, email string) error {
	if displayName == "" && email == "" {
		return ErrClientNoUpdate
	}
	return s.repo.UpdateUser(ctx, userID, displayName, email)
}
