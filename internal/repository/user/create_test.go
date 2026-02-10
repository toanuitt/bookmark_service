package user

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
				err := db.Where(whereByID, user.ID).First(&dbUser).Error
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
