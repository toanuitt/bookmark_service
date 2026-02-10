package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/toanuitt/bookmark_service/internal/handler/dto"
	"github.com/toanuitt/bookmark_service/internal/handler/utils"
	"github.com/toanuitt/bookmark_service/internal/service"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
	"github.com/toanuitt/bookmark_service/pkg/response"
)

const (
	// Error messages - Login
	invalidCredentialsError = "invalid username or password"
	loginSuccessMessage     = "Logged in successfully!"
)

// Login authenticates a user and returns a JWT token.
//
// @Summary Login a user
// @Description Authenticate a user with username and password, returns a JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param body body dto.LoginInputRequest true "User login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} response.Response "Invalid username or password"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/users/login [post]
func (h *user) Login(c *gin.Context) {
	input, err := utils.BindInputFromRequest[dto.LoginInputRequest](c)
	if err != nil {
		return
	}

	// call service
	token, err := h.svc.Login(c, input.Username, input.Password)
	switch {
	case errors.Is(err, service.ErrClientErr):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case errors.Is(err, dbutils.ErrNotFoundType):
		c.JSON(http.StatusNotFound, gin.H{"error": invalidCredentialsError})
		return
	case errors.Is(err, nil):
	default:
		log.Error().Err(err).Msg("Login: service returned error")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	// return token
	c.JSON(http.StatusOK, &dto.LoginResponse{
		Message: loginSuccessMessage,
		Data:    token,
	})
}
