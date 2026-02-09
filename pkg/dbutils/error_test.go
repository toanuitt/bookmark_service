package dbutils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCatchDBErr(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		inputErr    error
		expectedErr error
	}{
		{
			name:        "nil error",
			inputErr:    nil,
			expectedErr: nil,
		},
		{
			name:        "duplicate type error",
			inputErr:    errors.New("unique constraint failed: users.userID"),
			expectedErr: ErrDuplicationType,
		},
		{
			name:        "record not found error",
			inputErr:    gorm.ErrRecordNotFound,
			expectedErr: ErrNotFoundType,
		},
		{
			name:        "unknown error",
			inputErr:    errors.New("some random db error"),
			expectedErr: errors.New("some random db error"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := CatchDBErr(tc.inputErr)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
				return
			}

			assert.Equal(t, tc.expectedErr.Error(), err.Error())
		})
	}
}

func Test_filterDuplicationType(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		inputErr    error
		expectedHit bool
		expectedErr error
	}{
		{
			name:        "match unique constraint",
			inputErr:    errors.New("unique constraint failed: users.email"),
			expectedHit: true,
			expectedErr: ErrDuplicationType,
		},
		{
			name:        "not match",
			inputErr:    errors.New("some other error"),
			expectedHit: false,
			expectedErr: ErrDuplicationType,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hit, err := filterDuplicationType(tc.inputErr)

			assert.Equal(t, tc.expectedHit, hit)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func Test_filterRecordNotFound(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		inputErr    error
		expectedHit bool
		expectedErr error
	}{
		{
			name:        "gorm record not found",
			inputErr:    gorm.ErrRecordNotFound,
			expectedHit: true,
			expectedErr: ErrNotFoundType,
		},
		{
			name:        "other error",
			inputErr:    errors.New("some other error"),
			expectedHit: false,
			expectedErr: ErrNotFoundType,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hit, err := filterRecordNotFound(tc.inputErr)

			assert.Equal(t, tc.expectedHit, hit)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
