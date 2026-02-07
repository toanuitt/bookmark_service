package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/toanuitt/bookmark_service/internal/api"
	"github.com/toanuitt/bookmark_service/internal/model"
	jwtMocks "github.com/toanuitt/bookmark_service/pkg/jwtutils/mocks"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
	sqldbPkg "github.com/toanuitt/bookmark_service/pkg/sqldb"
	"github.com/toanuitt/bookmark_service/pkg/utils"
)

const (
	registerEndpoint = "/v1/users/register"
	loginEndpoint    = "/v1/users/login"
	selfInfoEndpoint = "/v1/self/info"

	headerContentType = "Content-Type"
	mimeJSON          = "application/json"
	headerAuth        = "Authorization"

	testUsername    = "John Doe"
	testBearerToken = "Bearer mock-token"
)

// TestUserRegisterEndpoint tests the user registration endpoint.
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

			app := api.New(cfg, redisPkg.InitMockRedis(t), db, jwtMocks.NewJWTGenerator(t), jwtMocks.NewJWTValidator(t))
			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestUserLoginEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                  string
		setupMockJWTGenerator func(t *testing.T) *jwtMocks.JWTGenerator
		setupTestHTTP         func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus        int
	}{
		{
			name: "success - valid user login",
			setupMockJWTGenerator: func(t *testing.T) *jwtMocks.JWTGenerator {
				jwtGen := jwtMocks.NewJWTGenerator(t)
				jwtGen.
					On("GenerateToken", mock.Anything).
					Return("mock-token", nil).
					Once()
				return jwtGen
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]string{
					"username": testUsername,
					"password": "TestPassword@123",
				}
				bodyBytes, err := json.Marshal(body)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, loginEndpoint, bytes.NewReader(bodyBytes))
				req.Header.Set(headerContentType, mimeJSON)

				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)

				if respRec.Code != http.StatusOK {
					t.Logf("Response body: %s", respRec.Body.String())
				}

				return respRec
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "bad request - invalid json",
			setupMockJWTGenerator: func(t *testing.T) *jwtMocks.JWTGenerator {
				return jwtMocks.NewJWTGenerator(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, loginEndpoint, bytes.NewReader([]byte(`{invalid json`)))
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
			jwtGen := tc.setupMockJWTGenerator(t)
			jwtVal := jwtMocks.NewJWTValidator(t)

			db := sqldbPkg.InitMockDb(t)
			err := db.AutoMigrate(&model.User{})
			require.NoError(t, err)

			hashedPwd := utils.HashPassword("TestPassword@123")
			err = db.Create(&model.User{
				Username: testUsername,
				Password: hashedPwd,
			}).Error
			require.NoError(t, err)

			app := api.New(
				cfg,
				redisPkg.InitMockRedis(t),
				db,
				jwtGen,
				jwtVal,
			)

			rec := tc.setupTestHTTP(app)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestGetProfileEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupJWT       func(t *testing.T) *jwtMocks.JWTValidator
		setupTestHTTP  func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name: "success - valid token",
			setupJWT: func(t *testing.T) *jwtMocks.JWTValidator {
				jwtVal := jwtMocks.NewJWTValidator(t)
				jwtVal.
					On("ValidateToken", mock.Anything).
					Return(jwt.MapClaims{
						"sub": "1",
					}, nil).
					Once()
				return jwtVal
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, selfInfoEndpoint, nil)
				req.Header.Set(headerAuth, testBearerToken)

				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "unauthorized - missing token",
			setupJWT: func(t *testing.T) *jwtMocks.JWTValidator {
				return jwtMocks.NewJWTValidator(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, selfInfoEndpoint, nil)

				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	cfg, err := api.NewConfig()
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwtGen := jwtMocks.NewJWTGenerator(t)
			jwtVal := tc.setupJWT(t)

			db := sqldbPkg.InitMockDb(t)
			err := db.AutoMigrate(&model.User{})
			require.NoError(t, err)

			err = db.Create(&model.User{
				ID:          "1",
				Username:    testUsername,
				DisplayName: "John",
				Email:       "john@example.com",
			}).Error
			require.NoError(t, err)

			app := api.New(
				cfg,
				redisPkg.InitMockRedis(t),
				db,
				jwtGen,
				jwtVal,
			)

			rec := tc.setupTestHTTP(app)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestUpdateProfileEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupJWT       func(t *testing.T) *jwtMocks.JWTValidator
		setupTestHTTP  func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name: "success - update display name",
			setupJWT: func(t *testing.T) *jwtMocks.JWTValidator {
				jwtVal := jwtMocks.NewJWTValidator(t)
				jwtVal.
					On("ValidateToken", mock.Anything).
					Return(jwt.MapClaims{
						"sub": "1",
					}, nil).
					Once()
				return jwtVal
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]string{
					"display_name": "New Name",
				}
				b, _ := json.Marshal(body)

				req := httptest.NewRequest(http.MethodPut, selfInfoEndpoint, bytes.NewReader(b))
				req.Header.Set(headerContentType, mimeJSON)
				req.Header.Set(headerAuth, testBearerToken)

				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "bad request - no fields",
			setupJWT: func(t *testing.T) *jwtMocks.JWTValidator {
				jwtVal := jwtMocks.NewJWTValidator(t)
				jwtVal.
					On("ValidateToken", mock.Anything).
					Return(jwt.MapClaims{
						"sub": "1",
					}, nil).
					Once()
				return jwtVal
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]string{}
				b, _ := json.Marshal(body)

				req := httptest.NewRequest(http.MethodPut, selfInfoEndpoint, bytes.NewReader(b))
				req.Header.Set(headerContentType, mimeJSON)
				req.Header.Set(headerAuth, testBearerToken)

				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "unauthorized - missing token",
			setupJWT: func(t *testing.T) *jwtMocks.JWTValidator {
				return jwtMocks.NewJWTValidator(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]string{
					"display_name": "New Name",
				}
				b, _ := json.Marshal(body)

				req := httptest.NewRequest(http.MethodPut, selfInfoEndpoint, bytes.NewReader(b))
				req.Header.Set(headerContentType, mimeJSON)

				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	cfg, err := api.NewConfig()
	require.NoError(t, err)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwtGen := jwtMocks.NewJWTGenerator(t)
			jwtVal := tc.setupJWT(t)

			db := sqldbPkg.InitMockDb(t)
			err := db.AutoMigrate(&model.User{})
			require.NoError(t, err)

			// Seed user
			err = db.Create(&model.User{
				ID:          "1",
				Username:    testUsername,
				DisplayName: "John",
				Email:       "john@example.com",
			}).Error
			require.NoError(t, err)

			app := api.New(
				cfg,
				redisPkg.InitMockRedis(t),
				db,
				jwtGen,
				jwtVal,
			)

			rec := tc.setupTestHTTP(app)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
