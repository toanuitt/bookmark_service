package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
)

func TestPasswordHanlder_GenPass(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func() *mocks.Password
		expectedStatus int
		expectedResp   string
	}{
		{
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
			},
			setupMockSvc: func() *mocks.Password {
				svcMock := mocks.NewPassword(t)
				svcMock.On("GeneratePassword").Return("123456789", nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedResp:   "123456789",
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
			},
			setupMockSvc: func() *mocks.Password {
				svcMock := mocks.NewPassword(t)
				svcMock.On("GeneratePassword").Return("", errors.New("something"))
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   "err",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc()
			testHanlder := NewPassword(mockSvc)
			testHanlder.GenPass(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResp, rec.Body.String())
		})
	}
}
