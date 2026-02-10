package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/internal/test/fixture"
	"gorm.io/gorm"
)

func TestUserRepo_UpdateUser(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		setupDB      func(t *testing.T) *gorm.DB
		inputUserID  string
		inputDisplay string
		inputEmail   string
		expectedErr  error
		verifyFunc   func(db *gorm.DB, userID string)
	}{
		{
			name: "normal case - update both display name and email",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.UserTestDB{})
				db.Create(&model.User{
					ID:          testUserID,
					DisplayName: testDisplayName,
					Username:    testUsername,
					Password:    testPasswordHash,
					Email:       testEmail,
				})
				return db
			},
			inputUserID:  testUserID,
			inputDisplay: testNewDisplayName,
			inputEmail:   testNewEmail,
			expectedErr:  nil,
			verifyFunc: func(db *gorm.DB, userID string) {
				var dbUser model.User
				err := db.Where(whereByID, userID).First(&dbUser).Error
				require.NoError(t, err)
				assert.Equal(t, testNewDisplayName, dbUser.DisplayName)
				assert.Equal(t, testNewEmail, dbUser.Email)
				assert.Equal(t, testUsername, dbUser.Username)
			},
		},
		{
			name: "update only display name",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.UserTestDB{})
				db.Create(&model.User{
					ID:          testUserID,
					DisplayName: testDisplayName,
					Username:    testUsername,
					Password:    testPasswordHash,
					Email:       testEmail,
				})
				return db
			},
			inputUserID:  testUserID,
			inputDisplay: testNewDisplayName,
			inputEmail:   "",
			expectedErr:  nil,
			verifyFunc: func(db *gorm.DB, userID string) {
				var dbUser model.User
				err := db.Where(whereByID, userID).First(&dbUser).Error
				require.NoError(t, err)
				assert.Equal(t, testNewDisplayName, dbUser.DisplayName)
				assert.Equal(t, testEmail, dbUser.Email)
			},
		},
		{
			name: "update only email",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.UserTestDB{})
				db.Create(&model.User{
					ID:          testUserID,
					DisplayName: testDisplayName,
					Username:    testUsername,
					Password:    testPasswordHash,
					Email:       testEmail,
				})
				return db
			},
			inputUserID:  testUserID,
			inputDisplay: "",
			inputEmail:   testNewEmail,
			expectedErr:  nil,
			verifyFunc: func(db *gorm.DB, userID string) {
				var dbUser model.User
				err := db.Where(whereByID, userID).First(&dbUser).Error
				require.NoError(t, err)
				assert.Equal(t, testDisplayName, dbUser.DisplayName)
				assert.Equal(t, testNewEmail, dbUser.Email)
			},
		},
		{
			name: "update with empty values - no changes",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.UserTestDB{})
				db.Create(&model.User{
					ID:          testUserID,
					DisplayName: testDisplayName,
					Username:    testUsername,
					Password:    testPasswordHash,
					Email:       testEmail,
				})
				return db
			},
			inputUserID:  testUserID,
			inputDisplay: "",
			inputEmail:   "",
			expectedErr:  nil,
			verifyFunc: func(db *gorm.DB, userID string) {
				var dbUser model.User
				err := db.Where(whereByID, userID).First(&dbUser).Error
				require.NoError(t, err)
				assert.Equal(t, testDisplayName, dbUser.DisplayName)
				assert.Equal(t, testEmail, dbUser.Email)
			},
		},
		{
			name: "user not found - UpdateUser",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserTestDB{})
			},
			inputUserID:  "019c134b-ca37-75bc-927a-881f2ce5c999",
			inputDisplay: testNewDisplayName,
			inputEmail:   testNewEmail,
			expectedErr:  nil,
			verifyFunc: func(db *gorm.DB, userID string) {
				var count int64
				db.Model(&model.User{}).Where(whereByID, userID).Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testRepo := NewUserRepository(db)

			err := testRepo.UpdateUser(ctx, tc.inputUserID, tc.inputDisplay, tc.inputEmail)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			require.NoError(t, err)

			if tc.verifyFunc != nil {
				tc.verifyFunc(db, tc.inputUserID)
			}
		})
	}
}
