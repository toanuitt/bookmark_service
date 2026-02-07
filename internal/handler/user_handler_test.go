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

// Test helpers

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

// testContext encapsulates common test context setup
type testContext struct {
	recorder *httptest.ResponseRecorder
	ginCtx   *gin.Context
}

func newTestContext() *testContext {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	gc, _ := gin.CreateTestContext(rec)
	return &testContext{
		recorder: rec,
		ginCtx:   gc,
	}
}

func (tc *testContext) setUserID(userID string) {
	tc.ginCtx.Set("userID", userID)
}

func (tc *testContext) assertStatusCode(t *testing.T, expected int) {
	assert.Equal(t, expected, tc.recorder.Code)
}

func (tc *testContext) assertBodyContains(t *testing.T, substring string) {
	assert.Contains(t, tc.recorder.Body.String(), substring)
}

func (tc *testContext) assertBodyNotEmpty(t *testing.T) {
	assert.NotEmpty(t, tc.recorder.Body.String())
}

func (tc *testContext) unmarshalResponse(t *testing.T, v any) {
	err := json.Unmarshal(tc.recorder.Body.Bytes(), v)
	assert.NoError(t, err)
}

// Mock service builders

func mockRegisterSuccess(t *testing.T, ctx *gin.Context) *mocks.Userservice {
	mockSvc := mocks.NewUserservice(t)
	mockSvc.On("Register", ctx, testUser, testPass, testName, testEmail).
		Return(&model.User{
			ID:          testID,
			Username:    testUser,
			DisplayName: testName,
			Email:       testEmail,
		}, nil)
	return mockSvc
}

func mockRegisterError(t *testing.T, err error) *mocks.Userservice {
	mockSvc := mocks.NewUserservice(t)
	mockSvc.On("Register", mock.Anything, testUser, testPass, testName, testEmail).
		Return(nil, err)
	return mockSvc
}

func mockLoginSuccess(t *testing.T, ctx *gin.Context) *mocks.Userservice {
	mockSvc := mocks.NewUserservice(t)
	mockSvc.On("Login", ctx, testUser, testPass).
		Return(testToken, nil)
	return mockSvc
}

func mockLoginError(t *testing.T, err error) *mocks.Userservice {
	mockSvc := mocks.NewUserservice(t)
	mockSvc.On("Login", mock.Anything, testUser, testPass).
		Return("", err)
	return mockSvc
}

func mockGetUserSuccess(t *testing.T, ctx *gin.Context) *mocks.Userservice {
	mockSvc := mocks.NewUserservice(t)
	mockSvc.On("GetUserByID", ctx, testID).
		Return(&model.User{
			ID:          testID,
			Username:    testUser,
			DisplayName: testName,
			Email:       testEmail,
		}, nil)
	return mockSvc
}

func mockGetUserError(t *testing.T, err error) *mocks.Userservice {
	mockSvc := mocks.NewUserservice(t)
	mockSvc.On("GetUserByID", mock.Anything, testID).
		Return(nil, err)
	return mockSvc
}

func mockUpdateUserSuccess(t *testing.T, ctx *gin.Context, displayName, email string) *mocks.Userservice {
	mockSvc := mocks.NewUserservice(t)
	mockSvc.On("UpdateUser", ctx, testID, displayName, email).
		Return(nil)
	return mockSvc
}

func mockUpdateUserError(t *testing.T, displayName, email string, err error) *mocks.Userservice {
	mockSvc := mocks.NewUserservice(t)
	mockSvc.On("UpdateUser", mock.Anything, testID, displayName, email).
		Return(err)
	return mockSvc
}

// Request builders

func buildRegisterRequest(t *testing.T) *http.Request {
	reqBody := &registerInputBody{
		Username:    testUser,
		Password:    testPass,
		DisplayName: testName,
		Email:       testEmail,
	}
	return NewJSONRequest(t, http.MethodPost, registerEndpoint, reqBody)
}

func buildLoginRequest(t *testing.T) *http.Request {
	reqBody := &loginInputBody{
		Username: testUser,
		Password: testPass,
	}
	return NewJSONRequest(t, http.MethodPost, loginEndpoint, reqBody)
}

func buildUpdateProfileRequest(t *testing.T, displayName, email string) *http.Request {
	reqBody := &updateProfileInputBody{
		DisplayName: displayName,
		Email:       email,
	}
	return NewJSONRequest(t, http.MethodPut, profileEndpoint, reqBody)
}

// Response validators

func assertRegisterResponse(t *testing.T, body string) {
	var res registerResBody
	err := json.Unmarshal([]byte(body), &res)
	assert.NoError(t, err)
	assert.Equal(t, "Register an user successfully!", res.Message)
	assert.NotNil(t, res.Data)
	assert.Equal(t, testID, res.Data.ID)
	assert.Equal(t, testUser, res.Data.Username)
	assert.Equal(t, testName, res.Data.DisplayName)
	assert.Equal(t, testEmail, res.Data.Email)
}

func assertLoginResponse(t *testing.T, body string) {
	var res loginResBody
	err := json.Unmarshal([]byte(body), &res)
	assert.NoError(t, err)
	assert.Equal(t, "Logged in successfully!", res.Message)
	assert.Equal(t, testToken, res.Data)
}

func assertProfileResponse(t *testing.T, body string) {
	var res profileResBody
	err := json.Unmarshal([]byte(body), &res)
	assert.NoError(t, err)
	assert.NotNil(t, res.Data)
	assert.Equal(t, testID, res.Data.ID)
	assert.Equal(t, testUser, res.Data.Username)
	assert.Equal(t, testName, res.Data.DisplayName)
	assert.Equal(t, testEmail, res.Data.Email)
}

// Test cases

func TestUserRegister(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		mockSvc        func(t *testing.T, ctx *gin.Context) *mocks.Userservice
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:           "success case",
			mockSvc:        mockRegisterSuccess,
			expectedStatus: http.StatusOK,
			checkResponse:  assertRegisterResponse,
		},
		{
			name: "duplicate username or email",
			mockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				return mockRegisterError(t, dbutils.ErrDuplicationType)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "username or email is already taken")
			},
		},
		{
			name: "service internal error",
			mockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				return mockRegisterError(t, errors.New("database error"))
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

			ctx := newTestContext()
			ctx.ginCtx.Request = buildRegisterRequest(t)
			mockSvc := tc.mockSvc(t, ctx.ginCtx)

			handler := NewUser(mockSvc)
			handler.RegisterUser(ctx.ginCtx)

			ctx.assertStatusCode(t, tc.expectedStatus)
			tc.checkResponse(t, ctx.recorder.Body.String())
		})
	}
}

func TestUserLogin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		mockSvc        func(t *testing.T, ctx *gin.Context) *mocks.Userservice
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:           "success case",
			mockSvc:        mockLoginSuccess,
			expectedStatus: http.StatusOK,
			checkResponse:  assertLoginResponse,
		},
		{
			name: "user not found",
			mockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				return mockLoginError(t, dbutils.ErrNotFoundType)
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "invalid username or password")
			},
		},
		{
			name: "client error - wrong password",
			mockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				return mockLoginError(t, service.ErrClientErr)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "error")
			},
		},
		{
			name: "internal server error",
			mockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				return mockLoginError(t, errors.New("database connection failed"))
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

			ctx := newTestContext()
			ctx.ginCtx.Request = buildLoginRequest(t)
			mockSvc := tc.mockSvc(t, ctx.ginCtx)

			handler := NewUser(mockSvc)
			handler.Login(ctx.ginCtx)

			ctx.assertStatusCode(t, tc.expectedStatus)
			tc.checkResponse(t, ctx.recorder.Body.String())
		})
	}
}

func TestUserGetProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		mockSvc        func(t *testing.T, ctx *gin.Context) *mocks.Userservice
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:           "success case",
			mockSvc:        mockGetUserSuccess,
			expectedStatus: http.StatusOK,
			checkResponse:  assertProfileResponse,
		},
		{
			name: "service returns error",
			mockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				return mockGetUserError(t, errors.New("database error"))
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

			ctx := newTestContext()
			ctx.ginCtx.Request = httptest.NewRequest(http.MethodGet, profileEndpoint, nil)
			ctx.setUserID(testID)
			mockSvc := tc.mockSvc(t, ctx.ginCtx)

			handler := NewUser(mockSvc)
			handler.GetProfile(ctx.ginCtx)

			ctx.assertStatusCode(t, tc.expectedStatus)
			tc.checkResponse(t, ctx.recorder.Body.String())
		})
	}
}

func TestUserUpdateProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		displayName    string
		email          string
		mockSvc        func(t *testing.T, ctx *gin.Context, displayName, email string) *mocks.Userservice
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:           "success case - update both fields",
			displayName:    testName,
			email:          testEmail,
			mockSvc:        mockUpdateUserSuccess,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, updateSelfInfoSuccessMessage)
			},
		},
		{
			name:           "success case - update display name only",
			displayName:    testName,
			email:          "",
			mockSvc:        mockUpdateUserSuccess,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, updateSelfInfoSuccessMessage)
			},
		},
		{
			name:           "success case - update email only",
			displayName:    "",
			email:          "newemail@example.com",
			mockSvc:        mockUpdateUserSuccess,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, updateSelfInfoSuccessMessage)
			},
		},
		{
			name:        "no fields provided for update",
			displayName: "",
			email:       "",
			mockSvc: func(t *testing.T, ctx *gin.Context, displayName, email string) *mocks.Userservice {
				return mockUpdateUserError(t, displayName, email, service.ErrClientNoUpdate)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "at least one field must be provided for update")
			},
		},
		{
			name:        "duplicate email error",
			displayName: "",
			email:       "existing@example.com",
			mockSvc: func(t *testing.T, ctx *gin.Context, displayName, email string) *mocks.Userservice {
				return mockUpdateUserError(t, displayName, email, dbutils.ErrDuplicationType)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "email is already taken")
			},
		},
		{
			name:        "internal server error",
			displayName: testName,
			email:       "",
			mockSvc: func(t *testing.T, ctx *gin.Context, displayName, email string) *mocks.Userservice {
				return mockUpdateUserError(t, displayName, email, errors.New("database connection failed"))
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

			ctx := newTestContext()
			ctx.ginCtx.Request = buildUpdateProfileRequest(t, tc.displayName, tc.email)
			ctx.setUserID(testID)
			mockSvc := tc.mockSvc(t, ctx.ginCtx, tc.displayName, tc.email)

			handler := NewUser(mockSvc)
			handler.UpdateProfile(ctx.ginCtx)

			ctx.assertStatusCode(t, tc.expectedStatus)
			tc.checkResponse(t, ctx.recorder.Body.String())
		})
	}
}
