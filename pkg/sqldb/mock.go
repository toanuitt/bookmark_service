package sqldb

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
