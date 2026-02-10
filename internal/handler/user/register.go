package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/toanuitt/bookmark_service/internal/handler/dto"
	"github.com/toanuitt/bookmark_service/internal/handler/utils"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
	"github.com/toanuitt/bookmark_service/pkg/response"
)

const (
	// Error messages - Registration
	usernameTakenError     = "username or email is already taken"
	registerSuccessMessage = "Register an user successfully!"
)

// RegisterUser handles user registration.
//
// @Summary Register a new user
// @Description Create a new user account using username, password, display name, and email.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration payload"
// @Success 200 {object} dto.RegisterResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/users/register [post]
func (h *user) RegisterUser(c *gin.Context) {
	input, err := utils.BindInputFromRequest[dto.RegisterRequest](c)
	if err != nil {
		return
	}

	res, err := h.svc.Register(c, input.Username, input.Password, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, dbutils.ErrDuplicationType):
		c.JSON(http.StatusBadRequest, gin.H{"error": usernameTakenError})
		return
	case errors.Is(err, nil):
	default:
		log.Error().Err(err).Msg("RegisterUser: service returned error")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, &dto.RegisterResponse{
		Message: registerSuccessMessage,
		Data: &dto.RegisterUserData{
			ID:          res.ID,
			Username:    res.Username,
			DisplayName: res.DisplayName,
			Email:       res.Email,
			UpdatedAt:   res.UpdatedAt.String(),
		},
	})
}
