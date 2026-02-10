package user

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/toanuitt/bookmark_service/internal/service"
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
)

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
