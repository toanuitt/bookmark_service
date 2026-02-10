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
	// Error messages - Profile Update
	noUpdateDataError            = "at least one field must be provided for update"
	emailTakenError              = "email is already taken"
	updateSelfInfoSuccessMessage = "Edit current user successfully!"
)

// UpdateProfile handles updating the current user's profile.
//
// @Summary Update current user profile
// @Description Update display name and/or email of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body dto.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} response.Response "Profile updated successfully"
// @Failure 400 {object} response.Response "No data provided for update"
// @Failure 401 {object} response.Response "Invalid or missing token"
// @Failure 404 {object} response.Response "User does not exist"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/self/info [put]
func (h *user) UpdateProfile(c *gin.Context) {
	userIDvalue, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidTokenError})
		return
	}

	userID, ok := userIDvalue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidTokenError})
		return
	}

	// Bind and validate input
	input, err := utils.BindInputFromRequest[dto.UpdateProfileRequest](c)
	if err != nil {
		return
	}

	// Call service to update user
	err = h.svc.UpdateUser(c, userID, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, service.ErrClientNoUpdate):
		c.JSON(http.StatusBadRequest, gin.H{"error": noUpdateDataError})
		return
	case errors.Is(err, dbutils.ErrDuplicationType):
		c.JSON(http.StatusBadRequest, gin.H{"error": emailTakenError})
		return
	case errors.Is(err, nil):
		c.JSON(http.StatusOK, &response.Response{
			Message: updateSelfInfoSuccessMessage,
		})
		return
	default:
		log.Error().Err(err).Str("user_id", userID).Msg("UpdateProfile: service error")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}
}
