package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/toanuitt/bookmark_service/internal/model"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepo --filename=user_repo.go

// UserRepo defines the interface for user persistence operations.
type UserRepo interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
}

var (
	// ErrDuplicateKey is returned when a unique constraint is violated (duplicate username).
	ErrDuplicateKey = errors.New("UNIQUE constraint failed: users.username")
)

type user struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepo instance.
func NewUserRepository(db *gorm.DB) UserRepo {
	return &user{db: db}
}

// CreateUser inserts a new user into the database.
// It returns ErrDuplicateKey if the username already exists.
func (r *user) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return nil, ErrDuplicateKey
		}
		return nil, err
	}
	return newUser, nil
}
