package user

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
)

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
