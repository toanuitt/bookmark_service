package dto

// HealthCheckResponse represents the JSON response body for the health check endpoint.
type HealthCheckResponse struct {
	Error       string `json:"error,omitempty"`
	Message     string `json:"message" example:"OK"`
	ServiceName string `json:"service_name" example:"bookmark_service"`
	InstanceID  string `json:"instance_id" example:"instance-test"`
}
