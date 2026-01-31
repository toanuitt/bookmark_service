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
				reqBody := map[string]string{
					"username":     "john_doe",
					"password":     "password123",
					"display_name": "John Doe",
					"email":        "john@example.com",
				}
				jsonBody, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBuffer(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("Register", ctx, "john_doe", "password123", "John Doe", "john@example.com").
					Return(&model.User{
						ID:          "019c134b-582c-7c27-a385-d1bb1dca44c5",
						Username:    "john_doe",
						DisplayName: "John Doe",
						Email:       "john@example.com",
					}, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var user model.User
				err := json.Unmarshal([]byte(body), &user)
				assert.NoError(t, err)
				assert.Equal(t, "john_doe", user.Username)
				assert.Equal(t, "John Doe", user.DisplayName)
				assert.Equal(t, "john@example.com", user.Email)
			},
		},
		{
			name: "invalid request body",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBuffer([]byte("invalid json")))
				ctx.Request.Header.Set("Content-Type", "application/json")
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
			name: "missing required fields",
			setupRequest: func(ctx *gin.Context) {
				reqBody := map[string]string{
					"username": "john_doe",
				}
				jsonBody, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBuffer(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
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
					"username":     "john_doe",
					"password":     "password123",
					"display_name": "John Doe",
					"email":        "john@example.com",
				}
				jsonBody, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBuffer(jsonBody))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.Userservice {
				mockSvc := mocks.NewUserservice(t)
				mockSvc.On("Register", mock.Anything, "john_doe", "password123", "John Doe", "john@example.com").
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
