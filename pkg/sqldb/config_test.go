package sqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	testcases := []struct {
		name     string
		env      map[string]string
		expected *config
	}{
		{
			name: "use default values when env not set",
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
			name: "override values from env",
			env: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_USER":     "myuser",
				"DB_PASSWORD": "mypassword",
				"DB_NAME":     "mydb",
				"DB_PORT":     "6543",
			},
			expected: &config{
				Host:     "db.example.com",
				User:     "myuser",
				Password: "mypassword",
				DBName:   "mydb",
				Port:     "6543",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			cfg, err := NewConfig("")
			require.NoError(t, err)
			require.NotNil(t, cfg)
			assert.Equal(t, tc.expected, cfg)
		})
	}
}

func TestConfig_GetDSN(t *testing.T) {
	t.Parallel()

	cfg := &config{
		Host:     "localhost",
		User:     "admin",
		Password: "secret",
		DBName:   "testdb",
		Port:     "5432",
	}

	expected := "host=localhost user=admin password=secret dbname=testdb port=5432 sslmode=disable"

	dsn := cfg.GetDSN()

	assert.Equal(t, expected, dsn)
}
