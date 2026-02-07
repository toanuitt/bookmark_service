package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/toanuitt/bookmark_service/internal/model"
	mockrepo "github.com/toanuitt/bookmark_service/internal/repository/mocks"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
	jwtMocks "github.com/toanuitt/bookmark_service/pkg/jwtutils/mocks"
	"github.com/toanuitt/bookmark_service/pkg/utils"
)

const (
	testUserID      = "user-123"
	testDisplayName = "Test User"
	testEmail       = "test@example.com"
	errUserNotFound = "user not found"
	updatedName     = "Updated Name"
	updatedEmail    = "updated@example.com"
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
			displayName: testDisplayName,
			email:       testEmail,
			setupMockRepo: func(t *testing.T, ctx context.Context, username, password, displayName, email string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)

				mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(user *model.User) bool {
					return user.Username == username &&
						user.DisplayName == displayName &&
						user.Email == email &&
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
				assert.Equal(t, testDisplayName, res.DisplayName)
				assert.Equal(t, testEmail, res.Email)
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
				})).Return(nil, dbutils.ErrDuplicationType)

				return mockRepo
			},
			expectedError: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockRepo := tc.setupMockRepo(t, ctx, tc.username, tc.password, tc.displayName, tc.email)
			jwtRepo := jwtMocks.NewJWTGenerator(t)
			service := NewUser(mockRepo, jwtRepo)

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

func TestUser_Login(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name           string
		username       string
		password       string
		setupMocks     func(t *testing.T, ctx context.Context, username, password string) (*mockrepo.UserRepo, *jwtMocks.JWTGenerator)
		expectedError  error
		validateResult func(t *testing.T, token string, err error)
	}{
		{
			name:     "successful login",
			username: "testuser",
			password: "password123",
			setupMocks: func(t *testing.T, ctx context.Context, username, password string) (*mockrepo.UserRepo, *jwtMocks.JWTGenerator) {
				mockRepo := mockrepo.NewUserRepo(t)
				mockJWT := jwtMocks.NewJWTGenerator(t)

				hashedPassword := utils.HashPassword(password)
				user := &model.User{
					ID:       testUserID,
					Username: username,
					Password: hashedPassword,
					Email:    testEmail,
				}

				mockRepo.On("GetUserByUsername", ctx, username).Return(user, nil)
				mockJWT.On("GenerateToken", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					sub, ok := claims["sub"].(string)
					return ok && sub == testUserID
				})).Return("test-jwt-token", nil)

				return mockRepo, mockJWT
			},
			expectedError: nil,
			validateResult: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "test-jwt-token", token)
			},
		},
		{
			name:     "user not found - login",
			username: "nonexistent",
			password: "password123",
			setupMocks: func(t *testing.T, ctx context.Context, username, password string) (*mockrepo.UserRepo, *jwtMocks.JWTGenerator) {
				mockRepo := mockrepo.NewUserRepo(t)
				mockJWT := jwtMocks.NewJWTGenerator(t)

				mockRepo.On("GetUserByUsername", ctx, username).Return(nil, errors.New(errUserNotFound))

				return mockRepo, mockJWT
			},
			expectedError: errors.New(errUserNotFound),
			validateResult: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Empty(t, token)
			},
		},
		{
			name:     "invalid password",
			username: "testuser",
			password: "wrongpassword",
			setupMocks: func(t *testing.T, ctx context.Context, username, password string) (*mockrepo.UserRepo, *jwtMocks.JWTGenerator) {
				mockRepo := mockrepo.NewUserRepo(t)
				mockJWT := jwtMocks.NewJWTGenerator(t)

				hashedPassword := utils.HashPassword("correctpassword")
				user := &model.User{
					ID:       testUserID,
					Username: username,
					Password: hashedPassword,
					Email:    testEmail,
				}

				mockRepo.On("GetUserByUsername", ctx, username).Return(user, nil)

				return mockRepo, mockJWT
			},
			expectedError: ErrClientErr,
			validateResult: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrClientErr, err)
				assert.Empty(t, token)
			},
		},
		{
			name:     "jwt generation error",
			username: "testuser",
			password: "password123",
			setupMocks: func(t *testing.T, ctx context.Context, username, password string) (*mockrepo.UserRepo, *jwtMocks.JWTGenerator) {
				mockRepo := mockrepo.NewUserRepo(t)
				mockJWT := jwtMocks.NewJWTGenerator(t)

				hashedPassword := utils.HashPassword(password)
				user := &model.User{
					ID:       testUserID,
					Username: username,
					Password: hashedPassword,
					Email:    testEmail,
				}

				mockRepo.On("GetUserByUsername", ctx, username).Return(user, nil)
				mockJWT.On("GenerateToken", mock.Anything).Return("", errors.New("jwt generation failed"))

				return mockRepo, mockJWT
			},
			expectedError: errors.New("jwt generation failed"),
			validateResult: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Empty(t, token)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockRepo, mockJWT := tc.setupMocks(t, ctx, tc.username, tc.password)
			service := NewUser(mockRepo, mockJWT)

			token, err := service.Login(ctx, tc.username, tc.password)

			tc.validateResult(t, token, err)
		})
	}
}

func TestUser_GetUserByID(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name          string
		userID        string
		setupMockRepo func(t *testing.T, ctx context.Context, userID string) *mockrepo.UserRepo
		expectedError error
		expectedUser  *model.User
	}{
		{
			name:   "successful get user by id",
			userID: testUserID,
			setupMockRepo: func(t *testing.T, ctx context.Context, userID string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)

				expectedUser := &model.User{
					ID:          userID,
					Username:    "testuser",
					Email:       testEmail,
					DisplayName: testDisplayName,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				mockRepo.On("GetUserById", ctx, userID).Return(expectedUser, nil)

				return mockRepo
			},
			expectedError: nil,
			expectedUser: &model.User{
				ID:          testUserID,
				Username:    "testuser",
				Email:       testEmail,
				DisplayName: testDisplayName,
			},
		},
		{
			name:   "user not found - getuserbyid",
			userID: "nonexistent-id",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)

				mockRepo.On("GetUserById", ctx, userID).Return(nil, errors.New(errUserNotFound))

				return mockRepo
			},
			expectedError: errors.New(errUserNotFound),
			expectedUser:  nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockRepo := tc.setupMockRepo(t, ctx, tc.userID)
			mockJWT := jwtMocks.NewJWTGenerator(t)
			service := NewUser(mockRepo, mockJWT)

			user, err := service.GetUserByID(ctx, tc.userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectedUser.ID, user.ID)
				assert.Equal(t, tc.expectedUser.Username, user.Username)
				assert.Equal(t, tc.expectedUser.Email, user.Email)
				assert.Equal(t, tc.expectedUser.DisplayName, user.DisplayName)
			}
		})
	}
}

func TestUser_UpdateUser(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name          string
		userID        string
		displayName   string
		email         string
		setupMockRepo func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockrepo.UserRepo
		expectedError error
	}{
		{
			name:        "successful update - both fields",
			userID:      testUserID,
			displayName: updatedName,
			email:       updatedEmail,
			setupMockRepo: func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)
				mockRepo.On("UpdateUser", ctx, userID, displayName, email).Return(nil)
				return mockRepo
			},
			expectedError: nil,
		},
		{
			name:        "successful update - display name only",
			userID:      testUserID,
			displayName: updatedName,
			email:       "",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)
				mockRepo.On("UpdateUser", ctx, userID, displayName, email).Return(nil)
				return mockRepo
			},
			expectedError: nil,
		},
		{
			name:        "successful update - email only",
			userID:      testUserID,
			displayName: "",
			email:       updatedEmail,
			setupMockRepo: func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)
				mockRepo.On("UpdateUser", ctx, userID, displayName, email).Return(nil)
				return mockRepo
			},
			expectedError: nil,
		},
		{
			name:        "no update - both fields empty",
			userID:      testUserID,
			displayName: "",
			email:       "",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)
				return mockRepo
			},
			expectedError: ErrClientNoUpdate,
		},
		{
			name:        "repository error",
			userID:      testUserID,
			displayName: updatedName,
			email:       updatedEmail,
			setupMockRepo: func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockrepo.UserRepo {
				mockRepo := mockrepo.NewUserRepo(t)
				mockRepo.On("UpdateUser", ctx, userID, displayName, email).Return(errors.New("database error"))
				return mockRepo
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockRepo := tc.setupMockRepo(t, ctx, tc.userID, tc.displayName, tc.email)
			mockJWT := jwtMocks.NewJWTGenerator(t)
			service := NewUser(mockRepo, mockJWT)

			err := service.UpdateUser(ctx, tc.userID, tc.displayName, tc.email)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
