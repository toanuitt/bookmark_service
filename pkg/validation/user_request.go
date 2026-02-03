package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	passwordRegex = regexp.MustCompile(`^.{8,}$`)
	upperRegex    = regexp.MustCompile(`[A-Z]`)
	lowerRegex    = regexp.MustCompile(`[a-z]`)
	numberRegex   = regexp.MustCompile(`[0-9]`)
	specialRegex  = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

// validateStrongPassword validates password using regex patterns:
// - At least 8 characters long
// - At least 1 uppercase letter
// - At least 1 lowercase letter
// - At least 1 number
// - At least 1 special character
func ValidateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return passwordRegex.MatchString(password) &&
		upperRegex.MatchString(password) &&
		lowerRegex.MatchString(password) &&
		numberRegex.MatchString(password) &&
		specialRegex.MatchString(password)
}
