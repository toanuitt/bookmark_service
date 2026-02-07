package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		inputErr      error
		expectedPanic bool
	}{
		{
			name:          "nil error - no panic",
			inputErr:      nil,
			expectedPanic: false,
		},
		{
			name:          "non-nil error - panic",
			inputErr:      errors.New("boom"),
			expectedPanic: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			didPanic := false

			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()

				HandleError(tc.inputErr)
			}()

			assert.Equal(t, tc.expectedPanic, didPanic)
		})
	}
}
