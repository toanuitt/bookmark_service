package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mockrepo "github.com/toanuitt/bookmark_service/internal/repository/mocks"
	mockcodegen "github.com/toanuitt/bookmark_service/pkg/stringutils/mocks"
)

var (
	dbErr       = errors.New("database connection failed")
	genErr      = errors.New("failed to generate code")
	storeErr    = errors.New("failed to store url")
	notFoundErr = errors.New("url not found")
)

func TestShortenURL_ShortenlengthURL(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		originURL   string
		expireAt    int64
		code        string
		setupMock   func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context)
		expectedRes string
		expectedErr error
	}{
		{
			name:      "normal case",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			code:      "abc1234",
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil)
				mockRepo.On("Exist", ctx, "abc1234").Return(false, nil)
				mockRepo.On("StoreUrl", ctx, "abc1234", "https://example.com", int64(3600)).Return(nil)
			},
			expectedRes: "abc1234",
			expectedErr: nil,
		},
		{
			name:      "mockRepo.Exist returns error",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil)
				mockRepo.On("Exist", ctx, "abc1234").Return(false, dbErr)
			},
			expectedRes: "",
			expectedErr: dbErr,
		},
		{
			name:      "max retry exceeded - code always exists",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil).Times(5)
				mockRepo.On("Exist", ctx, "abc1234").Return(true, nil).Times(5)
			},
			expectedRes: "",
			expectedErr: ErrMaxRetriesExceeded,
		},
		{
			name:      "code generator returns error",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context) {
				mockGen.On("GenerateCode", 7).Return("", genErr)
			},
			expectedRes: "",
			expectedErr: genErr,
		},
		{
			name:      "store url returns error",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil)
				mockRepo.On("Exist", ctx, "abc1234").Return(false, nil)
				mockRepo.On("StoreUrl", ctx, "abc1234", "https://example.com", int64(3600)).Return(storeErr)
			},
			expectedRes: "",
			expectedErr: storeErr,
		},
		{
			name:      "retry once then success",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil).Once()
				mockRepo.On("Exist", ctx, "abc1234").Return(true, nil).Once()
				mockGen.On("GenerateCode", 7).Return("xyz5678", nil).Once()
				mockRepo.On("Exist", ctx, "xyz5678").Return(false, nil).Once()
				mockRepo.On("StoreUrl", ctx, "xyz5678", "https://example.com", int64(3600)).Return(nil)
			},
			expectedRes: "xyz5678",
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gin.SetMode(gin.TestMode)

			mockRepo := mockrepo.NewUrlStorage(t)
			mockGen := mockcodegen.NewCodeGenerator(t)
			ctx := context.Background()
			tc.setupMock(mockRepo, mockGen, ctx)

			svc := NewShortenURL(mockRepo, mockGen)
			got, err := svc.ShortlengthURL(ctx, tc.originURL, tc.expireAt)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, got)
			}

			mockRepo.AssertExpectations(t)
			mockGen.AssertExpectations(t)
		})
	}
}

func TestShortenURL_GetURL(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		code        string
		setupMock   func(mockRepo *mockrepo.UrlStorage, ctx context.Context)
		expectedRes string
		expectedErr error
	}{
		{
			name: "normal case",
			code: "abc1234",
			setupMock: func(mockRepo *mockrepo.UrlStorage, ctx context.Context) {
				mockRepo.On("GetUrl", ctx, "abc1234").Return("https://example.com", nil)
			},
			expectedRes: "https://example.com",
			expectedErr: nil,
		},
		{
			name: "url not found",
			code: "invalid",
			setupMock: func(mockRepo *mockrepo.UrlStorage, ctx context.Context) {
				mockRepo.On("GetUrl", ctx, "invalid").Return("", notFoundErr)
			},
			expectedRes: "",
			expectedErr: notFoundErr,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gin.SetMode(gin.TestMode)

			mockRepo := mockrepo.NewUrlStorage(t)
			mockGen := mockcodegen.NewCodeGenerator(t)
			ctx := context.Background()
			tc.setupMock(mockRepo, ctx)

			svc := NewShortenURL(mockRepo, mockGen)
			got, err := svc.GetURL(ctx, tc.code)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
