package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mockrepo "github.com/toanuitt/bookmark_service/internal/repository/mocks"
	mockcodegen "github.com/toanuitt/bookmark_service/pkg/stringutils/mocks"
)

func TestShortenURL_ShortenlengthURL(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		originURL   string
		expireAt    int
		code        string
		setupMock   func(*mockrepo.UrlStorage, *mockcodegen.CodeGenerator)
		expectedRes string
		expectedErr bool
	}{
		{
			name:      "normal case",
			originURL: "https://example.com",
			expireAt:  3600,
			code:      "abc1234",
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil)
				mockRepo.On("Exist", mock.Anything, "abc1234").Return(false, nil)
				mockRepo.On("StoreUrl", mock.Anything, "abc1234", "https://example.com", 3600).Return(nil)
			},
			expectedRes: "abc1234",
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mockrepo.NewUrlStorage(t)
			mockGen := mockcodegen.NewCodeGenerator(t)
			tc.setupMock(mockRepo, mockGen)

			svc := NewShortenURL(mockRepo, mockGen)
			got, err := svc.ShortlengthURL(context.Background(), tc.originURL, tc.expireAt)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, got)
			}
		})
	}
}

func TestShortenURL_GetURL(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		code        string
		setupMock   func(*mockrepo.UrlStorage, *mockcodegen.CodeGenerator)
		expectedRes string
		expedtedErr bool
	}{
		{
			name: "normal case",
			code: "abc1234",
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator) {
				mockRepo.On("GetUrl", mock.Anything, "abc1234").Return("https://example.com", nil)
			},
			expectedRes: "https://example.com",
			expedtedErr: false,
		},
		{
			name: "url not found",
			code: "invalid",
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator) {
				mockRepo.On("GetUrl", mock.Anything, "invalid").Return("", errors.New("url not found"))
			},
			expectedRes: "",
			expedtedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mockrepo.NewUrlStorage(t)
			mockGen := mockcodegen.NewCodeGenerator(t)
			tc.setupMock(mockRepo, mockGen)

			svc := NewShortenURL(mockRepo, mockGen)
			got, err := svc.GetURL(context.Background(), tc.code)

			if tc.expedtedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, got)
			}
		})
	}
}
