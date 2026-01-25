package stringutils

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var urlSafeRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func TestCodeGeneratorPass_GenerateCode(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		length      int
		expectedLen int
		expectedErr error
	}{
		{
			name:        "normal case",
			length:      7,
			expectedLen: 7,
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testCgen := NewCodeGenerator()
			url, err := testCgen.GenerateCode(tc.length)
			assert.Equal(t, tc.expectedLen, len(url))
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, urlSafeRegex.MatchString(url), true)
		})
	}

}
