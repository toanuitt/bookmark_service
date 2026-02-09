package sqldb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected *config
	}{
		{
			name: "with defaults",
			env:  map[string]string{},
			expected: &config{
				Host:     "localhost",
				User:     "admin",
				Password: "admin",
				DBName:   "bookmark_service",
				Port:     "5432",
			},
		},
		{
			name: "with env overrides",
			env: map[string]string{
				"TEST_DB_HOST":     "db.example.com",
				"TEST_DB_USER":     "myuser",
				"TEST_DB_PASSWORD": "test-password",
				"TEST_DB_NAME":     "mydb",
				"TEST_DB_PORT":     "9999",
			},
			expected: &config{
				Host:     "db.example.com",
				User:     "myuser",
				Password: "test-password",
				DBName:   "mydb",
				Port:     "9999",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			cfg, err := NewConfig("TEST")

			require.NoError(t, err)
			require.NotNil(t, cfg)

			require.Equal(t, tc.expected.Host, cfg.Host)
			require.Equal(t, tc.expected.User, cfg.User)
			require.Equal(t, tc.expected.Password, cfg.Password)
			require.Equal(t, tc.expected.DBName, cfg.DBName)
			require.Equal(t, tc.expected.Port, cfg.Port)
		})
	}
}

func TestConfig_GetDSN(t *testing.T) {
	cfg := &config{
		Host:     "localhost",
		User:     "admin",
		Password: "test-pass",
		DBName:   "testdb",
		Port:     "5432",
	}

	dsn := cfg.GetDSN()

	expected := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
	)

	require.Equal(t, expected, dsn)
}
