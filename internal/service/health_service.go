package service

// healthService is the concrete implementation of the HealthCheck interface.
type healthService struct {
	serviceName string
	instanceID  string
}

//go:generate mockery --name HealthCheck --filename health_service.go

// HealthCheck defines the interface for health check operations.
type HealthCheck interface {
	CheckStatus() (string, string, string)
}

// NewHealthCheck creates and returns a new HealthCheck service instance.
// It takes the service name and instance ID as parameters.
func NewHealthCheck(serviceName string, instanceID string) HealthCheck {
	return &healthService{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

// CheckStatus returns the service health status, name, and instance ID.
// It always returns "OK" as the status message along with the configured service name and instance ID.
func (s *healthService) CheckStatus() (string, string, string) {
	return "OK", s.serviceName, s.instanceID
}
