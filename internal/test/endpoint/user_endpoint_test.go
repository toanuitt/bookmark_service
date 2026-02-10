package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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
	"gorm.io/gorm"
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
	testUserID      = "1"
	testPassword    = "TestPassword@123"
)

// =========================
// Test infrastructure
// =========================

func newTestApp(t *testing.T, jwtGen *jwtMocks.JWTGenerator, jwtVal *jwtMocks.JWTValidator) (api.Engine, *gorm.DB) {
	t.Helper()

	cfg, err := api.NewConfig()
	require.NoError(t, err)

	db := sqldbPkg.InitMockDb(t)
	require.NoError(t, db.AutoMigrate(&model.User{}))

	opts := &api.EngineOpts{
		Engine:       gin.New(),
		Cfg:          cfg,
		Redis:        redisPkg.InitMockRedis(t),
		SqlDB:        db,
		JWTGenerator: jwtGen,
		JWTValidator: jwtVal,
	}

	engine := api.New(opts)
	return engine, db
}

func seedUser(t *testing.T, db *gorm.DB, user *model.User) {
	t.Helper()
	err := db.Create(user).Error
	require.NoError(t, err)
}

func executeRequest(engine api.Engine, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

// =========================
// Request builders
// =========================

func newJSONRequest(t *testing.T, method, url string, body any) *http.Request {
	t.Helper()

	var bodyReader *bytes.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req := httptest.NewRequest(method, url, bodyReader)
	req.Header.Set(headerContentType, mimeJSON)
	return req
}

func newRawJSONRequest(method, url string, raw []byte) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewReader(raw))
	req.Header.Set(headerContentType, mimeJSON)
	return req
}

func newAuthenticatedRequest(t *testing.T, method, url string, body any, token string) *http.Request {
	t.Helper()
	req := newJSONRequest(t, method, url, body)
	req.Header.Set(headerAuth, token)
	return req
}

// =========================
// Mock builders
// =========================

func mockJWTGeneratorSuccess(t *testing.T) *jwtMocks.JWTGenerator {
	t.Helper()
	jwtGen := jwtMocks.NewJWTGenerator(t)
	jwtGen.On("GenerateToken", mock.Anything).
		Return("mock-token", nil).
		Once()
	return jwtGen
}

func mockJWTGeneratorNoop(t *testing.T) *jwtMocks.JWTGenerator {
	t.Helper()
	return jwtMocks.NewJWTGenerator(t)
}

func mockJWTValidatorSuccess(t *testing.T, userID string) *jwtMocks.JWTValidator {
	t.Helper()
	jwtVal := jwtMocks.NewJWTValidator(t)
	jwtVal.On("ValidateToken", mock.Anything).
		Return(jwt.MapClaims{"sub": userID}, nil).
		Once()
	return jwtVal
}

func mockJWTValidatorNoop(t *testing.T) *jwtMocks.JWTValidator {
	t.Helper()
	return jwtMocks.NewJWTValidator(t)
}

// =========================
// Test data builders
// =========================

func validRegisterBody() map[string]string {
	return map[string]string{
		"username":     "testuser",
		"password":     "SecurePass123!",
		"display_name": "Test User",
		"email":        "test@example.com",
	}
}

func registerBodyMissingField(field string) map[string]string {
	body := validRegisterBody()
	delete(body, field)
	return body
}

func validLoginBody() map[string]string {
	return map[string]string{
		"username": testUsername,
		"password": testPassword,
	}
}

func validUpdateProfileBody() map[string]string {
	return map[string]string{
		"display_name": "New Name",
	}
}

func defaultTestUser() *model.User {
	return &model.User{
		ID:          testUserID,
		Username:    testUsername,
		DisplayName: "John",
		Email:       "john@example.com",
		Password:    utils.HashPassword(testPassword),
	}
}

// =========================
// Tests
// =========================

func TestUserRegisterEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		requestBody    any
		rawBody        []byte
		expectedStatus int
	}{
		{
			name:           "success - valid user registration",
			requestBody:    validRegisterBody(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad request - invalid json",
			rawBody:        []byte("invalid json"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "bad request - missing username",
			requestBody:    registerBodyMissingField("username"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			engine, _ := newTestApp(t, mockJWTGeneratorNoop(t), mockJWTValidatorNoop(t))

			var req *http.Request
			if tc.rawBody != nil {
				req = newRawJSONRequest(http.MethodPost, registerEndpoint, tc.rawBody)
			} else {
				req = newJSONRequest(t, http.MethodPost, registerEndpoint, tc.requestBody)
			}

			rec := executeRequest(engine, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestUserLoginEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupJWT       func(t *testing.T) *jwtMocks.JWTGenerator
		requestBody    any
		rawBody        []byte
		seedUserFlag   bool
		expectedStatus int
	}{
		{
			name:           "success - valid user login",
			setupJWT:       mockJWTGeneratorSuccess,
			requestBody:    validLoginBody(),
			seedUserFlag:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad request - invalid json",
			setupJWT:       mockJWTGeneratorNoop,
			rawBody:        []byte(`{invalid json`),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwtGen := tc.setupJWT(t)
			engine, db := newTestApp(t, jwtGen, mockJWTValidatorNoop(t))

			if tc.seedUserFlag {
				seedUser(t, db, defaultTestUser())
			}

			var req *http.Request
			if tc.rawBody != nil {
				req = newRawJSONRequest(http.MethodPost, loginEndpoint, tc.rawBody)
			} else {
				req = newJSONRequest(t, http.MethodPost, loginEndpoint, tc.requestBody)
			}

			rec := executeRequest(engine, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			if rec.Code != http.StatusOK && tc.expectedStatus == http.StatusOK {
				t.Logf("Response body: %s", rec.Body.String())
			}
		})
	}
}

func TestGetProfileEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupJWT       func(t *testing.T) *jwtMocks.JWTValidator
		withAuthToken  bool
		expectedStatus int
	}{
		{
			name: "success - valid token",
			setupJWT: func(t *testing.T) *jwtMocks.JWTValidator {
				return mockJWTValidatorSuccess(t, testUserID)
			},
			withAuthToken:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "unauthorized - missing token",
			setupJWT:       mockJWTValidatorNoop,
			withAuthToken:  false,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwtVal := tc.setupJWT(t)
			engine, db := newTestApp(t, mockJWTGeneratorNoop(t), jwtVal)
			seedUser(t, db, defaultTestUser())

			req := httptest.NewRequest(http.MethodGet, selfInfoEndpoint, nil)
			if tc.withAuthToken {
				req.Header.Set(headerAuth, testBearerToken)
			}

			rec := executeRequest(engine, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestUpdateProfileEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupJWT       func(t *testing.T) *jwtMocks.JWTValidator
		requestBody    map[string]string
		withAuthToken  bool
		expectedStatus int
	}{
		{
			name: "success - update display name",
			setupJWT: func(t *testing.T) *jwtMocks.JWTValidator {
				return mockJWTValidatorSuccess(t, testUserID)
			},
			requestBody:    validUpdateProfileBody(),
			withAuthToken:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "bad request - no fields",
			setupJWT: func(t *testing.T) *jwtMocks.JWTValidator {
				return mockJWTValidatorSuccess(t, testUserID)
			},
			requestBody:    map[string]string{},
			withAuthToken:  true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unauthorized - missing token",
			setupJWT:       mockJWTValidatorNoop,
			requestBody:    validUpdateProfileBody(),
			withAuthToken:  false,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwtVal := tc.setupJWT(t)
			engine, db := newTestApp(t, mockJWTGeneratorNoop(t), jwtVal)
			seedUser(t, db, defaultTestUser())

			req := newJSONRequest(t, http.MethodPut, selfInfoEndpoint, tc.requestBody)
			if tc.withAuthToken {
				req.Header.Set(headerAuth, testBearerToken)
			}

			rec := executeRequest(engine, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
