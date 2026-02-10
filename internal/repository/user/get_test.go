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
