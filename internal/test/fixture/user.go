package fixture

import (
	"github.com/toanuitt/bookmark_service/internal/model"
	"gorm.io/gorm"
)

// UserTestDB is a test fixture helper for user-related database tests.
// It wraps a *gorm.DB and provides methods to:
//   - Setup the database connection
//   - Run migrations for the User model
//   - Seed test data for users
type UserTestDB struct {
	db *gorm.DB
}

// SetupDB assigns the given *gorm.DB to the fixture.
// This should be called before running migrations or generating test data.
func (f *UserTestDB) SetupDB(db *gorm.DB) {
	f.db = db
}

// Migrate runs GORM AutoMigrate for the User model.
// It ensures the users table (and related schema) exists before tests run.
func (f *UserTestDB) Migrate() error {
	return f.db.AutoMigrate(&model.User{})
}

// DB returns the underlying *gorm.DB instance.
// This is useful for tests that need direct access to the database.
func (f *UserTestDB) DB() *gorm.DB {
	return f.db
}

// GenerateData seeds the database with initial user records for testing.
// It inserts a batch of predefined users with fixed IDs, hashed passwords,
// and sample display names/emails to make tests deterministic.
func (f *UserTestDB) GenerateData() error {
	db := f.db.Session(&gorm.Session{})

	users := []*model.User{
		{
			ID:          "019c134b-582c-7c27-a385-d1bb1dca44c5",
			DisplayName: "John Doe",
			Username:    "John Doe",
			Password:    "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106",
			Email:       "john.doe@example.com",
		},
		{
			ID:          "019c134b-ca37-75bc-927a-881f2ce5c626",
			DisplayName: "Jane Doe",
			Username:    "Jane Doe",
			Password:    "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106",
			Email:       "jane.doe@example.com",
		},
	}

	// Insert users in batches to improve performance and keep the API consistent
	// with larger datasets in the future.
	return db.CreateInBatches(users, 10).Error
}
