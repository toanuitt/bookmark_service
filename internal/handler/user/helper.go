package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/toanuitt/bookmark_service/internal/handler/dto"
	"github.com/toanuitt/bookmark_service/internal/model"
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
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
	reqBody := &dto.RegisterRequest{
		Username:    testUser,
		Password:    testPass,
		DisplayName: testName,
		Email:       testEmail,
	}
	return NewJSONRequest(t, http.MethodPost, registerEndpoint, reqBody)
}

func buildLoginRequest(t *testing.T) *http.Request {
	reqBody := &dto.LoginInputRequest{
		Username: testUser,
		Password: testPass,
	}
	return NewJSONRequest(t, http.MethodPost, loginEndpoint, reqBody)
}

func buildUpdateProfileRequest(t *testing.T, displayName, email string) *http.Request {
	reqBody := &dto.UpdateProfileRequest{
		DisplayName: displayName,
		Email:       email,
	}
	return NewJSONRequest(t, http.MethodPut, profileEndpoint, reqBody)
}

// Response validators

func assertRegisterResponse(t *testing.T, body string) {
	var res dto.RegisterResponse
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
	var res dto.LoginResponse
	err := json.Unmarshal([]byte(body), &res)
	assert.NoError(t, err)
	assert.Equal(t, "Logged in successfully!", res.Message)
	assert.Equal(t, testToken, res.Data)
}

func assertProfileResponse(t *testing.T, body string) {
	var res dto.ProfileResponse
	err := json.Unmarshal([]byte(body), &res)
	assert.NoError(t, err)
	assert.NotNil(t, res.Data)
	assert.Equal(t, testID, res.Data.ID)
	assert.Equal(t, testUser, res.Data.Username)
	assert.Equal(t, testName, res.Data.DisplayName)
	assert.Equal(t, testEmail, res.Data.Email)
}
