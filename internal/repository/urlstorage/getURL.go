package urlstorage

import "context"

// GetUrl retrieves the original URL associated with the given shortened code.
func (r *urlStorage) GetUrl(ctx context.Context, code string) (string, error) {
	return r.c.Get(ctx, code).Result()
}
