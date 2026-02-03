package fixture

import (
	"testing"

	"github.com/toanuitt/bookmark_service/pkg/sqldb"
	"gorm.io/gorm"
)

// Fixture defines the contract for setting up a test database.
//
// A Fixture implementation is responsible for:
//   - Initializing the database connection (SetupDB)
//   - Migrating the schema (Migrate)
//   - Seeding test data (GenerateData)
//   - Exposing the underlying *gorm.DB (DB)
type Fixture interface {
	SetupDB(db *gorm.DB)
	Migrate() error
	GenerateData() error
	DB() *gorm.DB
}

// NewFixture initializes a test database using the given Fixture.
//
// It will:
//  1. Create a mock database and pass it to the fixture via SetupDB
//  2. Run schema migrations via Migrate
//  3. Generate test data via GenerateData
//
// If any step fails, the test will be aborted using t.Fatalf.
// On success, it returns the initialized *gorm.DB for use in tests.
func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	// create test database
	fix.SetupDB(sqldb.InitMockDb(t))

	// migrate schema
	err := fix.Migrate()
	if err != nil {
		t.Fatalf("Failed to migrate db for testing: %v", err)
	}

	// create test data
	err = fix.GenerateData()
	if err != nil {
		t.Fatalf("Failed to generate data for testing: %v", err)
	}

	return fix.DB()
}
