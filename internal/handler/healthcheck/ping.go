package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/toanuitt/bookmark_service/internal/handler/dto"
)

// @Summary Health check
// @Description Health check
// @Tags health_check
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse
// @Failure 500 {string} Internal Server Error
// @Router /v1/health-check [get]
func (h *healthHandler) CheckHealth(c *gin.Context) {
	message, serviceName, instanceID, err := h.svc.CheckStatus(c)
	if err != nil {
		log.Error().Err(err).Msg("Service return error on CheckHealth")
		c.JSON(http.StatusServiceUnavailable, dto.HealthCheckResponse{
			Error:       "Internal Server Error",
			Message:     "NOT OK",
			ServiceName: serviceName,
			InstanceID:  instanceID,
		})
		return
	}

	c.JSON(http.StatusOK, dto.HealthCheckResponse{
		Message:     message,
		ServiceName: serviceName,
		InstanceID:  instanceID,
	})
}
