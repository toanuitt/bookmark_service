package sqldb

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMockDb creates and initializes an in-memory SQLite database for testing.
//
// It uses a unique connection string (based on UUID) to ensure isolation between tests,
// while still allowing shared cache mode for multiple connections if needed.
// The GORM logger is set to Silent to keep test output clean.
//
// If the database cannot be initialized, the test is aborted using t.Fatal.
// On success, it returns a *gorm.DB instance ready to be used in tests.
func InitMockDb(t *testing.T) *gorm.DB {
	cxn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.Must(uuid.NewV7()).String())
	db, err := gorm.Open(sqlite.Open(cxn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal("Failed to initialize mock database:", err)
	}

	return db
}
