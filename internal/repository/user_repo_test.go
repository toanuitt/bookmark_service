package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/internal/test/fixture"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
	"gorm.io/gorm"
)

const (
	testUserID       = "019c134b-ca37-75bc-927a-881f2ce5c627"
	testDisplayName  = "John Doo"
	testUsername     = "John Doo"
	testPasswordHash = "$2a$10$wfpS7JvQgcHvHLk86eFs.jhKCIucgr9fhPkyBLVQntSHOnBOS106"
	testEmail        = "john.doo@example.com"

	testNewEmail       = "new@example.com"
	testNewUsername    = "Jane Doe"
	testNewDisplayName = "New User"
)

func TestUserRepo_CreateUser(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		setupDB        func(t *testing.T) *gorm.DB
		inputUser      *model.User
		expectedErr    error
		expectedOutput *model.User
		verifyFunc     func(db *gorm.DB, user *model.User)
	}{
		{
			name: "normal case",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserTestDB{})
			},
			inputUser: &model.User{
				ID:          testUserID,
				DisplayName: testDisplayName,
				Username:    testUsername,
				Password:    testPasswordHash,
				Email:       testEmail,
			},
			expectedErr: nil,
			expectedOutput: &model.User{
				ID:          testUserID,
				DisplayName: testDisplayName,
				Username:    testUsername,
				Password:    testPasswordHash,
				Email:       testEmail,
			},
			verifyFunc: func(db *gorm.DB, user *model.User) {
				var dbUser model.User
				err := db.Where("id = ?", user.ID).First(&dbUser).Error
				require.NoError(t, err)
				assert.Equal(t, user.Username, dbUser.Username)
				assert.Equal(t, user.Email, dbUser.Email)
				assert.Equal(t, user.DisplayName, dbUser.DisplayName)
			},
		},
		{
			name: "success-case create user without ID",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserTestDB{})
			},
			inputUser: &model.User{
				DisplayName: testDisplayName,
				Username:    testUsername,
				Password:    testPasswordHash,
				Email:       testEmail,
			},
			expectedErr:    nil,
			expectedOutput: nil,
			verifyFunc: func(db *gorm.DB, user *model.User) {
				assert.NotEmpty(t, user.ID)
				assert.Len(t, user.ID, 36)

				var dbUser model.User
				err := db.Where("username = ?", user.Username).First(&dbUser).Error
				require.NoError(t, err)
				assert.Equal(t, user.ID, dbUser.ID)
			},
		},
		{
			name: "duplicate username - fail",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserTestDB{})
			},
			inputUser: &model.User{
				ID:          testUserID,
				Username:    testNewUsername,
				Password:    testPasswordHash,
				Email:       testNewEmail,
				DisplayName: testNewDisplayName,
			},
			expectedErr:    dbutils.ErrDuplicationType,
			expectedOutput: nil,
			verifyFunc: func(db *gorm.DB, user *model.User) {
				var count int64
				db.Model(&model.User{}).Where("username = ?", "duplicate").Count(&count)
				assert.Equal(t, int64(1), count)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testRepo := NewUserRepository(db)

			result, err := testRepo.CreateUser(ctx, tc.inputUser)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			timeZero := time.Time{}
			result.CreatedAt = timeZero
			result.UpdatedAt = timeZero

			if tc.expectedOutput != nil {
				tc.expectedOutput.CreatedAt = timeZero
				tc.expectedOutput.UpdatedAt = timeZero
				assert.Equal(t, tc.expectedOutput, result)
			}

			if tc.verifyFunc != nil {
				tc.verifyFunc(db, result)
			}
		})
	}
}

func TestUserRepo_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		setupDB        func(t *testing.T) *gorm.DB
		inputUsername  string
		expectedErr    error
		expectedOutput *model.User
	}{
		{
			name: "normal case - user exists",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.UserTestDB{})
				// Create a test user
				db.Create(&model.User{
					ID:          testUserID,
					DisplayName: testDisplayName,
					Username:    testUsername,
					Password:    testPasswordHash,
					Email:       testEmail,
				})
				return db
			},
			inputUsername: testUsername,
			expectedErr:   nil,
			expectedOutput: &model.User{
				ID:          testUserID,
				DisplayName: testDisplayName,
				Username:    testUsername,
				Password:    testPasswordHash,
				Email:       testEmail,
			},
		},
		{
			name: "user not found - GetUserByUsername",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserTestDB{})
			},
			inputUsername:  "nonexistent",
			expectedErr:    dbutils.ErrNotFoundType,
			expectedOutput: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testRepo := NewUserRepository(db)

			result, err := testRepo.GetUserByUsername(ctx, tc.inputUsername)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			timeZero := time.Time{}
			result.CreatedAt = timeZero
			result.UpdatedAt = timeZero
			tc.expectedOutput.CreatedAt = timeZero
			tc.expectedOutput.UpdatedAt = timeZero

			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}

func TestUserRepo_GetUserById(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		setupDB        func(t *testing.T) *gorm.DB
		inputUserID    string
		expectedErr    error
		expectedOutput *model.User
	}{
		{
			name: "normal case - user exists",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.UserTestDB{})
				// Create a test user
				db.Create(&model.User{
					ID:          testUserID,
					DisplayName: testDisplayName,
					Username:    testUsername,
					Password:    testPasswordHash,
					Email:       testEmail,
				})
				return db
			},
			inputUserID: testUserID,
			expectedErr: nil,
			expectedOutput: &model.User{
				ID:          testUserID,
				DisplayName: testDisplayName,
				Username:    testUsername,
				Password:    testPasswordHash,
				Email:       testEmail,
			},
		},
		{
			name: "user not found - GetUserById",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserTestDB{})
			},
			inputUserID:    "019c134b-ca37-75bc-927a-881f2ce5c999",
			expectedErr:    dbutils.ErrNotFoundType,
			expectedOutput: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testRepo := NewUserRepository(db)

			result, err := testRepo.GetUserById(ctx, tc.inputUserID)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			timeZero := time.Time{}
			result.CreatedAt = timeZero
			result.UpdatedAt = timeZero
			tc.expectedOutput.CreatedAt = timeZero
			tc.expectedOutput.UpdatedAt = timeZero

			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}

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
				err := db.Where("id = ?", userID).First(&dbUser).Error
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
				err := db.Where("id = ?", userID).First(&dbUser).Error
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
				err := db.Where("id = ?", userID).First(&dbUser).Error
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
				err := db.Where("id = ?", userID).First(&dbUser).Error
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
				db.Model(&model.User{}).Where("id = ?", userID).Count(&count)
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
