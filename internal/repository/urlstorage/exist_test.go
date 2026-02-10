package urlstorage

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
)

func TestURLStorage_Exist(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name      string
		setupMock func() *redis.Client
		expected  bool
	}{
		{
			name: "key exists",
			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				mock.Set(context.Background(), "url", "https://google.com", 0)
				return mock
			},
			expected: true,
		},
		{
			name: "key not exists",
			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expected: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			redismock := tc.setupMock()
			testUrl := NewUrlStorage(redismock)
			exist, err := testUrl.Exist(ctx, "url")
			assert.Equal(t, tc.expected, exist)
			assert.NoError(t, err)
		})
	}
}
