package user

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
	"github.com/toanuitt/bookmark_service/pkg/dbutils"
)

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
