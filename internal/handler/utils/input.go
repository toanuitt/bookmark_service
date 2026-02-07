package utils

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/toanuitt/bookmark_service/pkg/response"
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

func BindInputFromRequest[T any](c *gin.Context) (*T, error) {
	reqInput := new(T)

	if err := c.ShouldBindJSON(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	if err := c.ShouldBindUri(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	if err := c.ShouldBindQuery(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	if err := c.ShouldBindHeader(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("password", ValidateStrongPassword)

	if err := validate.Struct(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	return reqInput, nil
}
