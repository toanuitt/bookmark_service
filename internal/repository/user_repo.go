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
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserById(ctx context.Context, userID string) (*model.User, error)
	UpdateUser(ctx context.Context, userID string, displayName string, email string) error
}

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
		return nil, dbutils.CatchDBErr(err)
	}
	return newUser, nil
}

func (r *user) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.GetUserByField(ctx, "username", username)
}

func (r *user) GetUserById(ctx context.Context, userID string) (*model.User, error) {
	return r.GetUserByField(ctx, "id", userID)
}

func (r *user) GetUserByField(ctx context.Context, field string, value string) (*model.User, error) {
	chosenUser := &model.User{}
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(chosenUser).Error
	if err != nil {
		return nil, dbutils.CatchDBErr(err)
	}

	return chosenUser, nil
}

func (r *user) UpdateUser(ctx context.Context, userID string, displayName string, email string) error {
	updates := make(map[string]any)
	if displayName != "" {
		updates["display_name"] = displayName
	}
	if email != "" {
		updates["email"] = email
	}
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return dbutils.CatchDBErr(err)
	}

	return nil
}
