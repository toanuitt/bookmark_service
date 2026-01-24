package service

import (
	"context"

	"github.com/toanuitt/bookmark_service/internal/repository"
)

// healthService is the concrete implementation of the HealthCheck interface.
type healthService struct {
	serviceName     string
	instanceID      string
	healthCheckRepo repository.HealthCheckRepo
}

//go:generate mockery --name HealthCheck --filename health_service.go

// HealthCheck defines the interface for health check operations.
type HealthCheck interface {
	CheckStatus(ctx context.Context) (string, string, string, error)
}

// NewHealthCheck creates and returns a new HealthCheck service instance.
// It takes the service name and instance ID as parameters.
func NewHealthCheck(serviceName string, instanceID string, healthCheckRepo repository.HealthCheckRepo) HealthCheck {
	return &healthService{
		serviceName:     serviceName,
		instanceID:      instanceID,
		healthCheckRepo: healthCheckRepo,
	}
}

// CheckStatus returns the service health status, name, and instance ID.
// It always returns "OK" as the status message along with the configured service name and instance ID.
func (s *healthService) CheckStatus(ctx context.Context) (string, string, string, error) {
	if err := s.healthCheckRepo.Ping(ctx); err != nil {
		return "redis: client is closed", s.serviceName, s.instanceID, err
	}
	return "OK", s.serviceName, s.instanceID, nil
}
