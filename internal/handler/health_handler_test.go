package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
)

func TestHealthCheckHanlder_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T) *mocks.HealthCheck
		expectedStatus int
		expectedResp   string
	}{
		{
			name: "normal case",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockSvc: func(t *testing.T) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("CheckStatus").Return("OK", "bookmark_service", "instance-test")
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedResp:   `{"message":"OK","service_name":"bookmark_service","instance_id":"instance-test"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(t)
			testHanlder := NewHealthCheck(mockSvc)
			testHanlder.CheckHealth(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResp, rec.Body.String())
		})
	}
}
