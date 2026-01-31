package fixture

import (
	"github.com/toanuitt/bookmark_service/internal/model"
	"gorm.io/gorm"
)

type UserTestDB struct {
	db *gorm.DB
}

func (f *UserTestDB) SetupDB(db *gorm.DB) {
	f.db = db

}

func (f *UserTestDB) Migrate() error {
	return f.db.AutoMigrate(&model.User{})
}

func (f *UserTestDB) DB() *gorm.DB {
	return f.db
}

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
	return db.CreateInBatches(users, 10).Error
}
