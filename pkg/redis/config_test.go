package redis

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newConfig(t *testing.T) {
	t.Parallel()

	oldEnv := map[string]string{
		"Redis_Addr":     os.Getenv("Redis_Addr"),
		"REDIS_PASSWORD": os.Getenv("REDIS_PASSWORD"),
		"REDIS_DB":       os.Getenv("REDIS_DB"),
	}

	t.Cleanup(func() {
		for k, v := range oldEnv {
			if v == "" {
				_ = os.Unsetenv(k)
			} else {
				_ = os.Setenv(k, v)
			}
		}
	})

	testcases := []struct {
		name     string
		env      map[string]string
		expected *config
	}{
		{
			name: "use default values when env not set",
			env:  map[string]string{},
			expected: &config{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			},
		},
		{
			name: "override values from env",
			env: map[string]string{
				"Redis_Addr":     "localhost:6379",
				"REDIS_PASSWORD": "secret",
				"REDIS_DB":       "2",
			},
			expected: &config{
				Addr:     "localhost:6379",
				Password: "secret",
				DB:       2,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for k, v := range tc.env {
				_ = os.Setenv(k, v)
			}
			cfg, err := newConfig("")
			require.NoError(t, err)
			require.NotNil(t, cfg)
			assert.Equal(t, tc.expected, cfg)
		})
	}
}
