package endpoint

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toanuitt/bookmark_service/internal/api"
	redisPkg "github.com/toanuitt/bookmark_service/pkg/redis"
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
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusOK,
			expectedCodeLen: 7,
			expectedMessage: "Shorten URL generated successfully!",
		},
		{
			name: "wrong input - empty url",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "",
					"exp": 10,
				}
				jsonBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusBadRequest,
			expectedCodeLen: 0,
			expectedMessage: "invalid request payload",
		},
	}

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := api.New(cfg, redisPkg.InitMockRedis(t))
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
