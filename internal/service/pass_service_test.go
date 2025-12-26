package service

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var urlSafeRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func TestPasswordService_GeneratePassword(t *testing.T) {
	t.Parallel()

	testcase := []struct {
		name        string
		expectedLen int
		expectErr   error
	}{
		{
			name:        "normal case",
			expectedLen: 10,
			expectErr:   nil,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewPassword()
			pass, err := testSvc.GeneratePassword()

			assert.Equal(t, tc.expectedLen, len(pass))
			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, urlSafeRegex.MatchString(pass), true)
		})
	}
}
