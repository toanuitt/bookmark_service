package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	mockrepo "github.com/toanuitt/bookmark_service/internal/repository/mocks"
)

func TestHealthCheckService_CheckStatus(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name                string
		serviceName         string
		instanceID          string
		setupMock           func(*mockrepo.HealthCheckRepo)
		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
		expectedErr         bool
	}{
		{
			name:        "success - redis is healthy",
			serviceName: "bookmark_service",
			instanceID:  "instance-test",
			setupMock: func(mockRepo *mockrepo.HealthCheckRepo) {
				mockRepo.On("Ping", context.Background()).Return(nil)
			},
			expectedMessage:     "OK",
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "instance-test",
			expectedErr:         false,
		},
		{
			name:        "redis connection error",
			serviceName: "bookmark_service",
			instanceID:  "instance-prod",
			setupMock: func(mockRepo *mockrepo.HealthCheckRepo) {
				mockRepo.On("Ping", context.Background()).Return(errors.New("redis: client is closed"))
			},
			expectedMessage:     "redis: client is closed",
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "instance-prod",
			expectedErr:         true,
		},
		{
			name:        "different service name",
			serviceName: "url_shortener_api",
			instanceID:  "instance-01",
			setupMock: func(mockRepo *mockrepo.HealthCheckRepo) {
				mockRepo.On("Ping", context.Background()).Return(nil)
			},
			expectedMessage:     "OK",
			expectedServiceName: "url_shortener_api",
			expectedInstanceID:  "instance-01",
			expectedErr:         false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mockrepo.NewHealthCheckRepo(t)
			tc.setupMock(mockRepo)

			testSvc := NewHealthCheck(tc.serviceName, tc.instanceID, mockRepo)

			// Call the method
			message, serviceName, instanceID, err := testSvc.CheckStatus(context.Background())

			assert.Equal(t, tc.expectedMessage, message)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
