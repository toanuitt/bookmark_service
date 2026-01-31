package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// config holds database connection configuration loaded from environment variables.
type config struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	User     string `default:"admin" envconfig:"DB_USER"`
	Password string `default:"admin" envconfig:"DB_PASSWORD"`
	DBName   string `default:"bookmark_service" envconfig:"DB_NAME"`
	Port     string `default:"5432" envconfig:"DB_PORT"`
}

// NewConfig loads database configuration using the given environment variable prefix.
// It reads values from environment variables and applies default values when not provided.
func NewConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetDSN builds and returns the Postgres DSN string from the configuration.
func (cfg *config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
	)
}
