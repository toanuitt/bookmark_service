package user

import (
	"context"

	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
)

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
