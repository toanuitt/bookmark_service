package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

// Config holds the application configuration parameters.
// All fields are read from environment variables with sensible defaults.
type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark-api" envconfig:"SERVICE_NAME"`
	InstanceID  string `default:"" envconfig:"INSTANCE_ID"`
}

// NewConfig creates and returns a new Config instance by reading environment variables.
// If INSTANCE_ID is not provided, a UUIDv7 is automatically generated.
// It returns an error if environment variable processing or UUID generation fails.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	if cfg.InstanceID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		cfg.InstanceID = id.String()
	}
	return cfg, nil
}
