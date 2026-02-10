package urlstorage

import "context"

// StoreUrl stores a URL with its shortened code in Redis with a 24-hour expiration.
func (r *urlStorage) StoreUrl(ctx context.Context, code, url string, exp int64) error {
	return r.c.Set(ctx, code, url, urlExpTime).Err()
}
