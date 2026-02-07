package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/internal/service"
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
)

const (
	registerEndpoint = "/v1/users/register"
	loginEndpoint    = "/v1/users/login"
	profileEndpoint  = "/v1/self/info"
	testEmail        = "john@example.com"
	testUser         = "john_doe"
	testPass         = "password123"
	testName         = "John Doe"
	testID           = "019c134b-582c-7c27-a385-d1bb1dca44c5"
	testToken        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test"
)

func NewJSONRequest(t *testing.T, method, url string, body any) *http.Request {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body to JSON: %v", err)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewRawJSONRequest(method, url string, raw string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(raw)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// TestUserRegister tests the RegisterUser handler
func TestUserRegister(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.Userservice
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "success case - register",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &registerInputBody{
					Username:    testUser,
					Password:    testPass,
					DisplayName: testName,
					Email:       testEmail,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPost, registerEndpoint, reqBody)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("Register", ctx, testUser, testPass, testName, testEmail).
					Return(&model.User{
						ID:          testID,
						Username:    testUser,
						DisplayName: testName,
						Email:       testEmail,
					}, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res registerResBody
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)

				assert.Equal(t, "Register an user successfully!", res.Message)
				assert.NotNil(t, res.Data)

				assert.Equal(t, testID, res.Data.ID)
				assert.Equal(t, testUser, res.Data.Username)
				assert.Equal(t, testName, res.Data.DisplayName)
				assert.Equal(t, testEmail, res.Data.Email)
			},
		},
		{
			name: "duplicate username or email",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &registerInputBody{
					Username:    testUser,
					Password:    testPass,
					DisplayName: testName,
					Email:       testEmail,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPost, registerEndpoint, reqBody)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("Register", mock.Anything, testUser, testPass, testName, testEmail).
					Return(nil, dbutils.ErrDuplicationType)
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "username or email is already taken")
			},
		},
		{
			name: "service internal error",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &registerInputBody{
					Username:    testUser,
					Password:    testPass,
					DisplayName: testName,
					Email:       testEmail,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPost, registerEndpoint, reqBody)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("Register", mock.Anything, testUser, testPass, testName, testEmail).
					Return(nil, errors.New("database error"))
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)

			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(t, gc)

			testHandler := NewUser(mockSvc)
			testHandler.RegisterUser(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			tc.checkResponse(t, rec.Body.String())
		})
	}
}

// TestUserLogin tests the Login handler
func TestUserLogin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.Userservice
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "success case - login",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &loginInputBody{
					Username: testUser,
					Password: testPass,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPost, loginEndpoint, reqBody)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("Login", ctx, testUser, testPass).
					Return(testToken, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res loginResBody
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.Equal(t, "Logged in successfully!", res.Message)
				assert.Equal(t, testToken, res.Data)
			},
		},
		{
			name: "user not found",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &loginInputBody{
					Username: testUser,
					Password: testPass,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPost, loginEndpoint, reqBody)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("Login", mock.Anything, testUser, testPass).
					Return("", dbutils.ErrNotFoundType)
				return mockSvc
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "invalid username or password")
			},
		},
		{
			name: "client error - wrong password",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &loginInputBody{
					Username: testUser,
					Password: testPass,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPost, loginEndpoint, reqBody)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("Login", mock.Anything, testUser, testPass).
					Return("", service.ErrClientErr)
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "error")
			},
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &loginInputBody{
					Username: testUser,
					Password: testPass,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPost, loginEndpoint, reqBody)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("Login", mock.Anything, testUser, testPass).
					Return("", errors.New("database connection failed"))
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)

			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(t, gc)

			testHandler := NewUser(mockSvc)
			testHandler.Login(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			tc.checkResponse(t, rec.Body.String())
		})
	}
}

// TestUserGetProfile tests the GetProfile handler
func TestUserGetProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupContext   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.Userservice
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "success case - getprofile",
			setupContext: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, profileEndpoint, nil)
				ctx.Set("userID", testID)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("GetUserByID", ctx, testID).
					Return(&model.User{
						ID:          testID,
						Username:    testUser,
						DisplayName: testName,
						Email:       testEmail,
					}, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res profileResBody
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)

				assert.NotNil(t, res.Data)
				assert.Equal(t, testID, res.Data.ID)
				assert.Equal(t, testUser, res.Data.Username)
				assert.Equal(t, testName, res.Data.DisplayName)
				assert.Equal(t, testEmail, res.Data.Email)
			},
		},
		{
			name: "service returns error",
			setupContext: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, profileEndpoint, nil)
				ctx.Set("userID", testID)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("GetUserByID", mock.Anything, testID).
					Return(nil, errors.New("database error"))
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)

			tc.setupContext(gc)
			mockSvc := tc.setupMockSvc(t, gc)

			testHandler := NewUser(mockSvc)
			testHandler.GetProfile(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			tc.checkResponse(t, rec.Body.String())
		})
	}
}

// TestUserUpdateProfile tests the UpdateProfile handler
func TestUserUpdateProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.Userservice
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "success case - update both fields",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &updateProfileInputBody{
					DisplayName: testName,
					Email:       testEmail,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPut, profileEndpoint, reqBody)
				ctx.Set("userID", testID)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("UpdateUser", ctx, testID, testName, testEmail).
					Return(nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, updateSelfInfoSuccessMessage)
			},
		},
		{
			name: "success case - update display name only",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &updateProfileInputBody{
					DisplayName: testName,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPut, profileEndpoint, reqBody)
				ctx.Set("userID", testID)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("UpdateUser", ctx, testID, testName, "").
					Return(nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, updateSelfInfoSuccessMessage)
			},
		},
		{
			name: "success case - update email only",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &updateProfileInputBody{
					Email: "newemail@example.com",
				}
				ctx.Request = NewJSONRequest(t, http.MethodPut, profileEndpoint, reqBody)
				ctx.Set("userID", testID)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("UpdateUser", ctx, testID, "", "newemail@example.com").
					Return(nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, updateSelfInfoSuccessMessage)
			},
		},
		{
			name: "no fields provided for update",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &updateProfileInputBody{}
				ctx.Request = NewJSONRequest(t, http.MethodPut, profileEndpoint, reqBody)
				ctx.Set("userID", testID)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("UpdateUser", ctx, testID, "", "").
					Return(service.ErrClientNoUpdate)
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "at least one field must be provided for update")
			},
		},
		{
			name: "duplicate email error",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &updateProfileInputBody{
					Email: "existing@example.com",
				}
				ctx.Request = NewJSONRequest(t, http.MethodPut, profileEndpoint, reqBody)
				ctx.Set("userID", testID)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("UpdateUser", ctx, testID, "", "existing@example.com").
					Return(dbutils.ErrDuplicationType)
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "email is already taken")
			},
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &updateProfileInputBody{
					DisplayName: testName,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPut, profileEndpoint, reqBody)
				ctx.Set("userID", testID)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("UpdateUser", ctx, testID, testName, "").
					Return(errors.New("database connection failed"))
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)

			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(t, gc)

			testHandler := NewUser(mockSvc)
			testHandler.UpdateProfile(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			tc.checkResponse(t, rec.Body.String())
		})
	}
}
