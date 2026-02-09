package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/toanuitt/bookmark_service/internal/api"
)

// TestPasswordEndpoint tests the password endpoint of the API.
// It tests the normal case where a GET request to /gen-pass returns a password of length 10.
// The test case checks if the returned status code and response length match the expected values.
// The test is run in parallel mode to ensure that the test runs quickly.
func TestPasswordEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus  int
		expectedRespLen int
	}{
		{
			name: "success",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/gen-pass", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusOK,
			expectedRespLen: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			app := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    &api.Config{},
			})
			rec := tc.setupTestHTTP(app)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedRespLen, rec.Body.Len())
		})
	}
}
