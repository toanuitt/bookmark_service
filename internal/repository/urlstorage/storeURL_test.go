package urlstorage

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
)

func TestUrlStorage_StoreUrl(t *testing.T) {
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
			name: "redis connection",
			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
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
			err := testUrl.StoreUrl(ctx, "url", "https://google.com", 3600)
			assert.Equal(t, tc.expectedErr, err)
			if err == nil {
				tc.verifyFunc(ctx, redismock)
			}
		})
	}
}
