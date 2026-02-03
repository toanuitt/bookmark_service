package sqldb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewClient creates and returns a new GORM database client using Postgres.
// It loads configuration using the given environment variable prefix,
// builds the DSN, and opens a connection to the database.
func NewClient(envPrefix string) (*gorm.DB, error) {
	cfg, err := NewConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
