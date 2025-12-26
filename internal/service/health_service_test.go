package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckService_CheckStatus(t *testing.T) {
	t.Parallel()
	testcase := []struct {
		name                string
		serviceName         string
		instanceId          string
		expectedMessage     string
		expectedServiceName string
		expectedInstanceId  string
	}{
		{
			name:                "normal case",
			serviceName:         "bookmark_service",
			instanceId:          "instance-test",
			expectedMessage:     "OK",
			expectedServiceName: "bookmark_service",
			expectedInstanceId:  "instance-test",
		},
	}
	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewHealthCheck(tc.serviceName, tc.instanceId)

			// Call the method
			message, serviceName, instanceID := testSvc.CheckStatus()

			assert.Equal(t, tc.expectedMessage, message)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceId, instanceID)
		})
	}
}
