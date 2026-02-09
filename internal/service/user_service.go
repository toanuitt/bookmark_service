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
	// tokenLast defines how long the JWT token is valid.
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
	// Register creates a new user with the given credentials and profile info.
	// It returns the created user or an error if creation fails.
	Register(ctx context.Context, username, password, displayName, email string) (*model.User, error)

	// Login verifies user credentials and returns a signed JWT token if successful.
	// It returns ErrClientErr if the username or password is invalid.
	Login(ctx context.Context, username, password string) (string, error)

	// GetUserByID retrieves a user by their unique ID.
	// It returns the user or an error if the user is not found or the query fails.
	GetUserByID(ctx context.Context, userId string) (*model.User, error)

	// UpdateUser updates the user's display name and/or email.
	// It returns ErrClientNoUpdate if both displayName and email are empty.
	UpdateUser(ctx context.Context, userID, displayName, email string) error
}

// user is the concrete implementation of Userservice.
type user struct {
	repo   repository.UserRepo
	jwtGen jwtutils.JWTGenerator
}

// NewUser creates and returns a new Userservice implementation.
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

// Login authenticates a user by username and password.
// If authentication succeeds, it generates and returns a JWT token.
// It returns ErrClientErr if the password is invalid.
func (s *user) Login(ctx context.Context, username, password string) (string, error) {
	// Check if user exists
	chosenUser, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	// Verify password
	check := utils.VerifyPassword(password, chosenUser.Password)
	if !check {
		return "", ErrClientErr
	}

	// Create JWT claims
	jwtContent := jwt.MapClaims{
		"sub": chosenUser.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tokenLast).Unix(),
	}

	// Generate token
	token, err := s.jwtGen.GenerateToken(jwtContent)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID returns a user by their ID.
func (s *user) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	return s.repo.GetUserById(ctx, userId)
}

// UpdateUser updates the user's display name and/or email.
// It returns ErrClientNoUpdate if no update fields are provided.
func (s *user) UpdateUser(ctx context.Context, userID, displayName, email string) error {
	if displayName == "" && email == "" {
		return ErrClientNoUpdate
	}
	return s.repo.UpdateUser(ctx, userID, displayName, email)
}
