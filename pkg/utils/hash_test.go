package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash_HashPassword(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		password string
	}{
		{
			name:     "normal case",
			password: "my-secret-password",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hash := HashPassword(tc.password)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tc.password, hash)
		})
	}
}

func TestHash_VerifyPassword(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		password    string
		input       string
		expectedRes bool
	}{
		{
			name:        "correct password",
			password:    "my-secret-password",
			input:       "my-secret-password",
			expectedRes: true,
		},
		{
			name:        "wrong password",
			password:    "my-secret-password",
			input:       "wrong-password",
			expectedRes: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hash := HashPassword(tc.password)
			ok := VerifyPassword(tc.input, hash)

			assert.Equal(t, tc.expectedRes, ok)
		})
	}
}
