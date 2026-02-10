package user

import (
	"context"

	"github.com/toanuitt/bookmark_service/internal/model"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepo --filename=user.go

// UserRepo defines the interface for user persistence operations.
type UserRepo interface {
	// CreateUser inserts a new user into the database.
	// It returns the created user or an error if the operation fails.
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)

	// GetUserByUsername retrieves a user by their username.
	// It returns the user or an error if the user is not found or the query fails.
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	// GetUserById retrieves a user by their ID.
	// It returns the user or an error if the user is not found or the query fails.
	GetUserById(ctx context.Context, userID string) (*model.User, error)

	// UpdateUser updates the user's display name and/or email by user ID.
	// It returns an error if the update operation fails.
	UpdateUser(ctx context.Context, userID string, displayName string, email string) error
}

// user is the GORM-based implementation of UserRepo.
type user struct {
	db *gorm.DB
}

// NewUserRepository creates and returns a new UserRepo backed by GORM.
func NewUserRepository(db *gorm.DB) UserRepo {
	return &user{db: db}
}
