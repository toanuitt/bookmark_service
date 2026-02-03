package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Response represents a standard API response structure.
type Response struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

var (
	// InternalErrorResponse is returned when an unexpected server error occurs.
	InternalErrorResponse = Response{Message: "Something went wrong", Details: nil}

	// InvalidRequestError is returned when the request payload is invalid.
	InvalidRequestError = Response{Message: "Invalid request", Details: nil}
)

// InputFieldError converts validation errors into a user-friendly response.
// If the error is not a validator.ValidationErrors type, it returns InternalErrorResponse instead.
func InputFieldError(e error) Response {
	if ok := errors.As(e, &validator.ValidationErrors{}); !ok {
		return InternalErrorResponse
	}

	var errs []string
	for _, err := range e.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid "+err.Tag())
	}

	return Response{Message: "Invalid request", Details: errs}
}
