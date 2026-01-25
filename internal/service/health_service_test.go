package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	mockrepo "github.com/toanuitt/bookmark_service/internal/repository/mocks"
)

var (
	TestRedisClosedErr = errors.New("redis: client is closed")
)

func TestHealthCheckService_CheckStatus(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name                string
		serviceName         string
		instanceID          string
		setupMock           func(mockRepo *mockrepo.HealthCheckRepo, ctx context.Context)
		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
		expectedErr         error
	}{
		{
			name:        "success - redis is healthy",
			serviceName: "bookmark_service",
			instanceID:  "instance-test",
			setupMock: func(mockRepo *mockrepo.HealthCheckRepo, ctx context.Context) {
				mockRepo.On("Ping", ctx).Return(nil)
			},
			expectedMessage:     "OK",
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "instance-test",
			expectedErr:         nil,
		},
		{
			name:        "redis connection error",
			serviceName: "bookmark_service",
			instanceID:  "instance-prod",
			setupMock: func(mockRepo *mockrepo.HealthCheckRepo, ctx context.Context) {
				mockRepo.On("Ping", ctx).Return(errors.New("redis: client is closed"))
			},
			expectedMessage:     "redis: client is closed",
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "instance-prod",
			expectedErr:         TestRedisClosedErr,
		},
		{
			name:        "different service name",
			serviceName: "url_shortener_api",
			instanceID:  "instance-01",
			setupMock: func(mockRepo *mockrepo.HealthCheckRepo, ctx context.Context) {
				mockRepo.On("Ping", ctx).Return(nil)
			},
			expectedMessage:     "OK",
			expectedServiceName: "url_shortener_api",
			expectedInstanceID:  "instance-01",
			expectedErr:         nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mockrepo.NewHealthCheckRepo(t)
			ctx := context.Background()
			tc.setupMock(mockRepo, ctx)

			testSvc := NewHealthCheck(tc.serviceName, tc.instanceID, mockRepo)

			// Call the method
			message, serviceName, instanceID, err := testSvc.CheckStatus(ctx)

			assert.Equal(t, tc.expectedMessage, message)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
