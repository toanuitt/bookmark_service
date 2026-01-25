package redis

import "github.com/kelseyhightower/envconfig"

// config holds Redis connection configuration.
type config struct {
	Addr     string `default:"localhost:6379" envconfig:"Redis_Addr"`
	Password string `default:"" envconfig:"REDIS_PASSWORD"`
	DB       int    `default:"0" envconfig:"REDIS_DB"`
}

// newConfig creates and returns a new config instance by reading environment variables.
func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
