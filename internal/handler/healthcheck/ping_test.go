package healthcheck

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
)

var testConnectError = errors.New("can't connect to redis")

func TestHealthCheckHanlder_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.HealthCheck
		expectedStatus int
		expectedResp   string
	}{
		{
			name: "normal case",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("CheckStatus", ctx).Return("OK", "bookmark_service", "instance-test", nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedResp:   `{"message":"OK","service_name":"bookmark_service","instance_id":"instance-test"}`,
		},
		{
			name: "error case",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("CheckStatus", ctx).Return("OK", "bookmark_service", "instance-test", testConnectError)
				return mockSvc
			},

			expectedStatus: http.StatusServiceUnavailable,
			expectedResp:   `{"error":"Internal Server Error","message":"NOT OK","service_name":"bookmark_service","instance_id":"instance-test"}`,
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
			testHanlder := NewHealthCheck(mockSvc)
			testHanlder.CheckHealth(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResp, rec.Body.String())
		})
	}
}
