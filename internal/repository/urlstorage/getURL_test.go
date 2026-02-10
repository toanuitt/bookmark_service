package urlstorage

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
)

func TestURLStorage_GetURL(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		setupMock   func() *redis.Client
		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",
			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				mock.Set(context.Background(), "url", "https://google.com", 100)
				return mock
			},
			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				url, err := r.Get(ctx, "url").Result()
				assert.Nil(t, err)
				assert.Equal(t, url, "https://google.com")
			},
		},
		{
			name: "code not found",
			setupMock: func() *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectedErr: redis.Nil,
		},
		{
			name: "error redis closed",
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
			redismock := tc.setupMock()
			testUrl := NewUrlStorage(redismock)
			_, err := testUrl.GetUrl(ctx, "url")
			assert.Equal(t, tc.expectedErr, err)
			if err == nil {
				tc.verifyFunc(ctx, redismock)
			}
		})
	}
}
