package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	t.Parallel()

	oldEnv := os.Getenv("LOG_LEVEL")
	oldLevel := zerolog.GlobalLevel()

	t.Cleanup(func() {
		_ = os.Setenv("LOG_LEVEL", oldEnv)
		zerolog.SetGlobalLevel(oldLevel)
	})

	testcases := []struct {
		name          string
		envValue      string
		expectedLevel zerolog.Level
	}{
		{
			name:          "valid level - debug",
			envValue:      "debug",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			name:          "valid level - warn",
			envValue:      "warn",
			expectedLevel: zerolog.WarnLevel,
		},
		{
			name:          "invalid level - fallback to info",
			envValue:      "not-a-level",
			expectedLevel: zerolog.InfoLevel,
		},
		{
			name:          "empty env - fallback to info",
			envValue:      "",
			expectedLevel: zerolog.InfoLevel,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_ = os.Setenv("LOG_LEVEL", tc.envValue)

			SetLogLevel()

			actual := zerolog.GlobalLevel()
			assert.Equal(t, tc.expectedLevel, actual)
		})
	}
}
