package repository

import (
	"context"
	"fmt"

	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepo --filename=user_repo.go

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

// CreateUser inserts a new user into the database.
// It returns ErrDuplicateKey (wrapped by dbutils) if the username already exists.
func (r *user) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, dbutils.CatchDBErr(err)
	}
	return newUser, nil
}

// GetUserByUsername retrieves a user by the "username" field.
func (r *user) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.GetUserByField(ctx, "username", username)
}

// GetUserById retrieves a user by the "id" field.
func (r *user) GetUserById(ctx context.Context, userID string) (*model.User, error) {
	return r.GetUserByField(ctx, "id", userID)
}

// GetUserByField retrieves a user by an arbitrary field and value.
// It returns an error if no record is found or the query fails.
func (r *user) GetUserByField(ctx context.Context, field string, value string) (*model.User, error) {
	chosenUser := &model.User{}
	err := r.db.WithContext(ctx).
		Where(fmt.Sprintf("%s = ?", field), value).
		First(chosenUser).Error
	if err != nil {
		return nil, dbutils.CatchDBErr(err)
	}

	return chosenUser, nil
}

// UpdateUser updates the user's display name and/or email by user ID.
// Only non-empty fields will be updated.
func (r *user) UpdateUser(ctx context.Context, userID string, displayName string, email string) error {
	updates := make(map[string]any)
	if displayName != "" {
		updates["display_name"] = displayName
	}
	if email != "" {
		updates["email"] = email
	}

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
	if err != nil {
		return dbutils.CatchDBErr(err)
	}

	return nil
}
