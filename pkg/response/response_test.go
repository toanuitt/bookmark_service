package response

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string `validate:"required"`
	Age  int    `validate:"gte=18"`
}

func TestInputFieldError(t *testing.T) {
	t.Parallel()

	validate := validator.New()

	testcases := []struct {
		name        string
		inputErr    error
		expectedRes Response
	}{
		{
			name:     "non validation error -> internal error",
			inputErr: errors.New("some random error"),
			expectedRes: Response{
				Message: InternalErrorResponse.Message,
				Details: InternalErrorResponse.Details,
			},
		},
		{
			name: "validation error - single field",
			inputErr: func() error {
				s := testStruct{
					Name: "",
					Age:  20,
				}
				return validate.Struct(s)
			}(),
			expectedRes: Response{
				Message: "Invalid request",
				Details: []string{
					"Name is invalid required",
				},
			},
		},
		{
			name: "validation error - multiple fields",
			inputErr: func() error {
				s := testStruct{
					Name: "",
					Age:  10,
				}
				return validate.Struct(s)
			}(),
			expectedRes: Response{
				Message: "Invalid request",
				Details: []string{
					"Name is invalid required",
					"Age is invalid gte",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := InputFieldError(tc.inputErr)

			assert.Equal(t, tc.expectedRes.Message, res.Message)
			assert.Equal(t, tc.expectedRes.Details, res.Details)
		})
	}
}
