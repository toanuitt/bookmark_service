package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toanuitt/bookmark_service/internal/service"
)

// healthCheckResponse represents the JSON response body for the health check endpoint.
type healthCheckResponse struct {
	Message     string `json:"message" example:"OK"`
	ServiceName string `json:"service_name" example:"bookmark_service"`
	InstanceID  string `json:"instance_id" example:"instance-test"`
}

// healthHandler is the concrete implementation of the HealthCheck handler interface.
type healthHandler struct {
	svc service.HealthCheck
}

// HealthCheck defines the interface for health check HTTP handlers.
type HealthCheck interface {
	CheckHealth(c *gin.Context)
}

// NewHealthCheck creates and returns a new HealthCheck handler instance.
// It takes a service.HealthCheck dependency to perform health status checks.
func NewHealthCheck(svc service.HealthCheck) HealthCheck {
	return &healthHandler{svc: svc}
}

// @Summary Health check
// @Description Health check
// @Tags health_check
// @Produce json
// @Success 200 {object} healthCheckResponse
// @Failure 500 {string} Internal Server Error
// @Router /health-check [get]
func (h *healthHandler) CheckHealth(c *gin.Context) {
	message, serviceName, instanceID := h.svc.CheckStatus()
	c.JSON(http.StatusOK, healthCheckResponse{
		Message:     message,
		ServiceName: serviceName,
		InstanceID:  instanceID,
	})
}
