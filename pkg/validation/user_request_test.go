package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Password string `validate:"strong_password"`
}

func TestValidateStrongPassword(t *testing.T) {
	t.Parallel()

	v := validator.New()
	err := v.RegisterValidation("strong_password", ValidateStrongPassword)
	assert.NoError(t, err)

	testcases := []struct {
		name       string
		password   string
		expectedOK bool
	}{
		{
			name:       "valid strong password",
			password:   "Abcdef1!",
			expectedOK: true,
		},
		{
			name:       "too short",
			password:   "Ab1!",
			expectedOK: false,
		},
		{
			name:       "missing uppercase",
			password:   "abcdef1!",
			expectedOK: false,
		},
		{
			name:       "missing lowercase",
			password:   "ABCDEF1!",
			expectedOK: false,
		},
		{
			name:       "missing number",
			password:   "Abcdefg!",
			expectedOK: false,
		},
		{
			name:       "missing special char",
			password:   "Abcdefg1",
			expectedOK: false,
		},
		{
			name:       "empty password",
			password:   "",
			expectedOK: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data := testStruct{
				Password: tc.password,
			}
			err := v.Struct(data)
			if tc.expectedOK {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
