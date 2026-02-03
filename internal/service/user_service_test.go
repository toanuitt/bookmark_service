package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/internal/repository"
	mockrepo "github.com/toanuitt/bookmark_service/internal/repository/mocks"
	"github.com/toanuitt/bookmark_service/pkg/utils"
)

func TestUser_Register(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name           string
		username       string
		password       string
		displayName    string
		email          string
		setupMockRepo  func(t *testing.T, ctx context.Context, username, password, displayName, email string) *mockrepo.UserRepo
		expectedError  error
		validateResult func(t *testing.T, res *model.User, err error)
	}{
		{
			name:        "successful registration",
			username:    "testuser",
			password:    "password123",
			displayName: "Test User",
			email:       "test@example.com",
			setupMockRepo: func(t *testing.T, ctx context.Context, username, password, displayName, email string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)
				hashedPassword := utils.HashPassword(password)

				expectedUser := &model.User{
					Username:    username,
					Password:    hashedPassword,
					DisplayName: displayName,
					Email:       email,
				}

				mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(user *model.User) bool {
					return user.Username == expectedUser.Username &&
						user.DisplayName == expectedUser.DisplayName &&
						user.Email == expectedUser.Email &&
						utils.VerifyPassword(password, user.Password) &&
						!user.CreatedAt.IsZero() &&
						!user.UpdatedAt.IsZero()
				})).Return(func(ctx context.Context, user *model.User) *model.User {
					user.ID = "test-uuid-123"
					return user
				}, nil)

				return mockRepo
			},
			expectedError: nil,
			validateResult: func(t *testing.T, res *model.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.ID)
				assert.Equal(t, "testuser", res.Username)
				assert.Equal(t, "Test User", res.DisplayName)
				assert.Equal(t, "test@example.com", res.Email)
				assert.NotEmpty(t, res.Password)
				assert.True(t, utils.VerifyPassword("password123", res.Password))
				assert.NotZero(t, res.CreatedAt)
				assert.NotZero(t, res.UpdatedAt)
			},
		},
		{
			name:        "repository error - user already exists",
			username:    "existinguser",
			password:    "password123",
			displayName: "Existing User",
			email:       "existing@example.com",
			setupMockRepo: func(t *testing.T, ctx context.Context, username, password, displayName, email string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)

				mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(user *model.User) bool {
					return user.Username == username &&
						user.DisplayName == displayName &&
						user.Email == email &&
						utils.VerifyPassword(password, user.Password)
				})).Return(nil, repository.ErrDuplicateKey)

				return mockRepo
			},
			expectedError: repository.ErrDuplicateKey,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRepo := tc.setupMockRepo(t, ctx, tc.username, tc.password, tc.displayName, tc.email)
			service := NewUser(mockRepo)

			resp, err := service.Register(ctx, tc.username, tc.password, tc.displayName, tc.email)

			if tc.validateResult != nil {
				tc.validateResult(t, resp, err)
			} else {
				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tc.expectedError, err)
					assert.Nil(t, resp)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, resp)
				}
			}
		})
	}
}
