package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toanuitt/bookmark_service/internal/service/mocks"
)

func TestShortenURL_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T, ctx *gin.Context) *mocks.ShortenURLservice
		expectedStatus int
		expectedResp   ShortenURLResponse
	}{
		{
			name: "normal case",
			setupRequest: func(ctx *gin.Context) {
				reqBody := ShortenURLRequest{
					URL:      "https://example.com",
					ExpireIn: 3600,
				}
				bodyBytes, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.ShortenURLservice {

				mockSvc := mocks.NewShortenURLservice(t)
				mockSvc.On("ShortlengthURL", ctx, "https://example.com", int64(3600)).Return("abc1234", nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedResp: ShortenURLResponse{
				Message: "Shorten URL generated successfully!",
				Code:    "abc1234",
			},
		},
		{
			name: "service error during URL shortening",
			setupRequest: func(ctx *gin.Context) {
				reqBody := ShortenURLRequest{
					URL:      "https://example.com",
					ExpireIn: 3600,
				}
				bodyBytes, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.ShortenURLservice {

				mockSvc := mocks.NewShortenURLservice(t)
				mockSvc.On("ShortlengthURL", ctx, "https://example.com", int64(3600)).Return("", assert.AnError)
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp: ShortenURLResponse{
				Message: "internal server error",
				Code:    "",
			},
		},
		{
			name: "invalid request payload",
			setupRequest: func(ctx *gin.Context) {
				reqBody := ShortenURLRequest{
					URL:      "https://example.com",
					ExpireIn: -1,
				}
				bodyBytes, _ := json.Marshal(reqBody)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			setupMockSvc: func(t *testing.T, ctx *gin.Context) *mocks.ShortenURLservice {
				mockSvc := mocks.NewShortenURLservice(t)
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: ShortenURLResponse{
				Message: "invalid request payload",
				Code:    "",
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
			testHandler := NewShortenURL(mockSvc)
			testHandler.ShortenURL(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			var resp ShortenURLResponse
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			require.NoError(t, err, "Failed to unmarshal response")
			assert.Equal(t, tc.expectedResp, resp)
			mockSvc.AssertExpectations(t)
		})
	}
}
