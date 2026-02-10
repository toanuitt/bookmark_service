package ping

import "context"

// Ping checks the health of the Redis connection.
func (r *healthCheckStorage) Ping(ctx context.Context) error {
	return r.redisClient.Ping(ctx).Err()
}
