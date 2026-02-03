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
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
)

const (
	endpoint  = "/v1/user/register"
	testEmail = "john@example.com"
	testUser  = "john_doe"
	testPass  = "password123"
	testName  = "John Doe"
	testID    = "019c134b-582c-7c27-a385-d1bb1dca44c5"
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

func TestUser_RegisterUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.Userservice
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "success case",
			setupRequest: func(ctx *gin.Context) {
				reqBody := &registerInputBody{
					Username:    testUser,
					Password:    testPass,
					DisplayName: testName,
					Email:       testEmail,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPost, endpoint, reqBody)
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
				var user model.User
				err := json.Unmarshal([]byte(body), &user)
				assert.NoError(t, err)
				assert.Equal(t, testUser, user.Username)
				assert.Equal(t, testName, user.DisplayName)
				assert.Equal(t, testEmail, user.Email)
			},
		},
		{
			name: "invalid request body",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = NewRawJSONRequest(http.MethodPost, endpoint, "invalid json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.NotEmpty(t, body)
			},
		},
		{
			name: "service error",
			setupRequest: func(ctx *gin.Context) {
				reqBody := map[string]string{
					"username":     testUser,
					"password":     testPass,
					"display_name": testName,
					"email":        testEmail,
				}
				ctx.Request = NewJSONRequest(t, http.MethodPost, endpoint, reqBody)
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
