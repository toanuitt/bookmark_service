package endpoint

import (
	"bytes"
	"context"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toanuitt/bookmark_service/internal/api"
	jwtMocks "github.com/toanuitt/bookmark_service/pkg/jwtutils/mocks"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
	sqldbPkg "github.com/toanuitt/bookmark_service/pkg/sqldb"
)

func TestUrlShortenEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus  int
		expectedCodeLen int
		expectedMessage string
	}{
		{
			name: "success",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "https://google.com",
					"exp": 10,
				}
				jsonBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusOK,
			expectedCodeLen: 7,
			expectedMessage: "Shorten URL generated successfully!",
		},
		{
			name: "invalid json",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader([]byte(`{invalid-json`)))
				req.Header.Set("Content-Type", "application/json")
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusBadRequest,
			expectedCodeLen: 0,
			expectedMessage: "Something went wrong",
		},
	}

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := api.New(cfg, redisPkg.InitMockRedis(t), sqldbPkg.InitMockDb(t), jwtMocks.NewJWTGenerator(t), jwtMocks.NewJWTValidator(t))
			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]string
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedMessage, resp["message"])
			if tc.expectedCodeLen > 0 {
				assert.Equal(t, tc.expectedCodeLen, len(resp["code"]))
			}

		})
	}
}
func TestUrlRedirectEndpoint(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		setupMock     func(ctx context.Context) *redis.Client
		setupTestHTTP func(app api.Engine) *httptest.ResponseRecorder

		expectedStatus int
		expectedHeader string
		expectedBody   string
	}{
		{
			name: "success redirect",

			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				err := mock.Set(ctx, "1234567", "https://google.com", 3000).Err()
				require.NoError(t, err)
				return mock
			},

			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodGet,
					"/v1/links/redirect/1234567",
					nil,
				)
				rec := httptest.NewRecorder()
				app.ServeHTTP(rec, req)
				return rec
			},

			expectedStatus: http.StatusFound,
			expectedHeader: "https://google.com",
		},
		{
			name: "not found",

			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},

			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodGet,
					"/v1/links/redirect/notfound",
					nil,
				)
				rec := httptest.NewRecorder()
				app.ServeHTTP(rec, req)
				return rec
			},

			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"url not found"}`,
		},
	}

	cfg, err := api.NewConfig()
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			redisMock := tc.setupMock(ctx)
			app := api.New(cfg, redisMock, sqldbPkg.InitMockDb(t), jwtMocks.NewJWTGenerator(t), jwtMocks.NewJWTValidator(t))

			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedHeader != "" {
				assert.Equal(t, tc.expectedHeader, rec.Header().Get("Location"))
				return
			}

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
		})
	}
}
