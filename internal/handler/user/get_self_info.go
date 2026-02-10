package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/toanuitt/bookmark_service/internal/handler/dto"
	"github.com/toanuitt/bookmark_service/pkg/response"
)

const (
	// Error messages - Authentication
	invalidTokenError = "Invalid Token"
)

// GetProfile returns the current user's profile based on JWT token.
//
// @Summary Get current user profile
// @Description Get profile of the currently authenticated user using JWT in Authorization header
// @Tags Users
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Success 200 {object} dto.ProfileResponse
// @Failure 401 {object} response.Response "Invalid or missing token"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/self/info [get]
func (h *user) GetProfile(c *gin.Context) {
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

	res, err := h.svc.GetUserByID(c, userID)
	if err != nil {
		log.Error().Err(err).Msg("Getprofile: service returned error")
		c.JSON(http.StatusInternalServerError, response.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, &dto.ProfileResponse{
		Data: res,
	})
}
