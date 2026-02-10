package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/toanuitt/bookmark_service/internal/service"
)

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
