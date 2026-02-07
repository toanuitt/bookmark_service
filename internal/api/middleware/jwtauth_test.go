package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	jwtMocks "github.com/toanuitt/bookmark_service/pkg/jwtutils/mocks"
)

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		setupMock      func(*jwtMocks.JWTValidator)
		expectedStatus int
		expectedBody   map[string]any
		expectClaims   bool
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"error": "Authorization header is required",
			},
			expectClaims: false,
		},
		{
			name:           "wrong authorization format",
			authHeader:     "BadFormatToken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"error": "Authorization header format is wrong",
			},
			expectClaims: false,
		},
		{
			name:       "invalid token - validator returns error",
			authHeader: "Bearer invalidtoken",
			setupMock: func(m *jwtMocks.JWTValidator) {
				m.On("ValidateToken", "invalidtoken").
					Return(jwt.MapClaims(nil), errors.New("invalid token")).
					Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"error": "Invalid Token",
			},
			expectClaims: false,
		},
		{
			name:       "valid token but missing sub claim",
			authHeader: "Bearer validtoken",
			setupMock: func(m *jwtMocks.JWTValidator) {
				m.On("ValidateToken", "validtoken").
					Return(jwt.MapClaims{
						"foo": "bar",
					}, nil).
					Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"error": "Invalid Token",
			},
			expectClaims: false,
		},
		{
			name:       "valid token with sub claim",
			authHeader: "Bearer validtoken",
			setupMock: func(m *jwtMocks.JWTValidator) {
				m.On("ValidateToken", "validtoken").
					Return(jwt.MapClaims{
						"sub": "user-123",
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"ok": true,
			},
			expectClaims: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()

			mockValidator := new(jwtMocks.JWTValidator)
			if tc.setupMock != nil {
				tc.setupMock(mockValidator)
			}

			jwtAuth := NewJWTAuth(mockValidator)
			router.Use(jwtAuth.JWTAuth())
			router.GET("/test", func(c *gin.Context) {
				if tc.expectClaims {
					uid, exists := c.Get("userID")
					assert.True(t, exists)
					assert.NotNil(t, uid)
				}
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
			var body map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &body)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, body)
			mockValidator.AssertExpectations(t)
		})
	}
}
