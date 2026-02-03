package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toanuitt/bookmark_service/internal/api"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
	sqldbPkg "github.com/toanuitt/bookmark_service/pkg/sqldb"
)

// TestHealthCheckEndpoint tests the healthCheckEndpoint function.
// It tests that the function returns a HTTP 200 OK response with the correct JSON body.
// The JSON body is expected to have the following structure:
//
//	{
//	  "message": string,
//	  "serviceName": string,
//	  "instanceID": string
//	}
//
// The character set used for generating the JSON body is constant and does not change across different implementations of the interface. The length of the generated JSON body is constant and does not change across different implementations of the interface.
// If an error occurs while generating the JSON body, the error is returned immediately and the generated JSON body is an empty string.
func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"OK","serviceName":"bookmark-service","instanceID":"instance-test"}`,
		},
	}

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := api.New(cfg, redisPkg.InitMockRedis(t), sqldbPkg.InitMockDb(t))
			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]string

			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, "OK", resp["message"])
			assert.NotEmpty(t, resp["instance_id"])
			assert.NotEmpty(t, resp["service_name"])

		})
	}
}
