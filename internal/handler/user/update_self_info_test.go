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
