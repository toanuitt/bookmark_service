package user

import (
	"context"
	"fmt"

	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
)

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
