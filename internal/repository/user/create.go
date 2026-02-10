package user

import (
	"context"

	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
)

// CreateUser inserts a new user into the database.
// It returns ErrDuplicateKey (wrapped by dbutils) if the username already exists.
func (r *user) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, dbutils.CatchDBErr(err)
	}
	return newUser, nil
}
