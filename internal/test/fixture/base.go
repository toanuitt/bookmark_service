package fixture

import (
	"testing"

	"github.com/toanuitt/bookmark_service/pkg/sqldb"
	"gorm.io/gorm"
)

type Fixture interface {
	SetupDB(db *gorm.DB)
	Migrate() error
	GenerateData() error
	DB() *gorm.DB
}

func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	//create test database
	fix.SetupDB(sqldb.InitMockDb(t))
	//migrate schema
	err := fix.Migrate()

	if err != nil {
		t.Fatalf("Failed to migrate db for testing: %v", err)
	}

	//create test data
	err = fix.GenerateData()
	if err != nil {
		t.Fatalf("Failed to generate data for testing: %v", err)
	}

	return fix.DB()
}
