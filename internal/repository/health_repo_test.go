package repository

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
)

func TestHealthRepo_Ping(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		setupMock   func() *redis.Client
		expectedErr error
	}{
		{
			name: "normal case",
			setupMock: func() *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			redisMock := tc.setupMock()
			healthCheck := NewHealthCheckStorage(redisMock)
			err := healthCheck.Ping(ctx)
			assert.Equal(t, tc.expectedErr, err)

		})
	}
}
