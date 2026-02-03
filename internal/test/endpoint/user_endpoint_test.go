package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toanuitt/bookmark_service/internal/api"
	"github.com/toanuitt/bookmark_service/internal/model"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
	sqldbPkg "github.com/toanuitt/bookmark_service/pkg/sqldb"
)

const (
	registerEndpoint  = "/v1/users/register"
	headerContentType = "Content-Type"
	mimeJSON          = "application/json"
)

// TestUserRegisterEndpoint tests the user registration endpoint.
// It tests that the function returns appropriate HTTP responses for valid and invalid inputs.
// The test covers the following scenarios:
// - Successful user registration with valid input (HTTP 200 OK)
// - Invalid request body (HTTP 400 Bad Request)
// - Missing required fields (HTTP 400 Bad Request)
func TestUserRegisterEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus int
	}{
		{
			name: "success - valid user registration",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]string{
					"username":     "testuser",
					"password":     "SecurePass123!",
					"display_name": "Test User",
					"email":        "test@example.com",
				}
				bodyBytes, err := json.Marshal(body)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, registerEndpoint, bytes.NewReader(bodyBytes))
				req.Header.Set(headerContentType, mimeJSON)

				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus: http.StatusOK,
		},
		{
			name: "bad request - invalid json",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := []byte("invalid json")

				req := httptest.NewRequest(http.MethodPost, registerEndpoint, bytes.NewReader(body))
				req.Header.Set(headerContentType, mimeJSON)

				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad request - missing username",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]string{
					"password":     "SecurePass123!",
					"display_name": "Test User",
					"email":        "test@example.com",
				}
				bodyBytes, err := json.Marshal(body)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, registerEndpoint, bytes.NewReader(bodyBytes))
				req.Header.Set(headerContentType, mimeJSON)

				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus: http.StatusBadRequest,
		},
	}

	cfg, err := api.NewConfig()
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := sqldbPkg.InitMockDb(t)
			require.NoError(t, db.AutoMigrate(&model.User{}))

			app := api.New(cfg, redisPkg.InitMockRedis(t), db)
			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
