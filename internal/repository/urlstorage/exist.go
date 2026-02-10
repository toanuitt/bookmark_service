package urlstorage

import "context"

// Exist checks whether a shortened URL code exists in Redis.
func (r *urlStorage) Exist(ctx context.Context, code string) (bool, error) {
	result, err := r.c.Exists(ctx, code).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}
