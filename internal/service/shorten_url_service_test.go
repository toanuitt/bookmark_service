package service

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mockrepo "github.com/toanuitt/bookmark_service/internal/repository/mocks"
	mockcodegen "github.com/toanuitt/bookmark_service/pkg/stringutils/mocks"
)

func TestShortenURL_ShortenlengthURL(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		originURL   string
		expireAt    int64
		code        string
		setupMock   func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx *gin.Context) // Sửa: context.Context -> *gin.Context
		expectedRes string
		expectedErr bool
	}{
		{
			name:      "normal case",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			code:      "abc1234",
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx *gin.Context) { // Sửa: context.Context -> *gin.Context
				mockGen.On("GenerateCode", 7).Return("abc1234", nil)
				mockRepo.On("Exist", ctx, "abc1234").Return(false, nil)
				mockRepo.On("StoreUrl", ctx, "abc1234", "https://example.com", int64(3600)).Return(nil)
			},
			expectedRes: "abc1234",
			expectedErr: false,
		},
		{
			name:      "mockRepo.Exist returns error",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx *gin.Context) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil)
				mockRepo.On("Exist", ctx, "abc1234").Return(false, errors.New("database connection failed"))
			},
			expectedRes: "",
			expectedErr: true,
		},
		{
			name:      "max retry exceeded - code always exists",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx *gin.Context) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil)
				mockRepo.On("Exist", ctx, "abc1234").Return(true, nil).Times(5)
			},
			expectedRes: "",
			expectedErr: true,
		},
		{
			name:      "code generator returns error",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx *gin.Context) {
				mockGen.On("GenerateCode", 7).Return("", errors.New("failed to generate code"))
			},
			expectedRes: "",
			expectedErr: true,
		},
		{
			name:      "store url returns error",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx *gin.Context) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil)
				mockRepo.On("Exist", ctx, "abc1234").Return(false, nil)
				mockRepo.On("StoreUrl", ctx, "abc1234", "https://example.com", int64(3600)).Return(errors.New("failed to store url"))
			},
			expectedRes: "",
			expectedErr: true,
		},
		{
			name:      "retry once then success",
			originURL: "https://example.com",
			expireAt:  int64(3600),
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx *gin.Context) {
				mockGen.On("GenerateCode", 7).Return("abc1234", nil).Once()
				mockRepo.On("Exist", ctx, "abc1234").Return(true, nil).Once()
				mockGen.On("GenerateCode", 7).Return("xyz5678", nil).Once()
				mockRepo.On("Exist", ctx, "xyz5678").Return(false, nil).Once()
				mockRepo.On("StoreUrl", ctx, "xyz5678", "https://example.com", int64(3600)).Return(nil)
			},
			expectedRes: "xyz5678",
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gin.SetMode(gin.TestMode)

			mockRepo := mockrepo.NewUrlStorage(t)
			mockGen := mockcodegen.NewCodeGenerator(t)
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupMock(mockRepo, mockGen, gc)

			svc := NewShortenURL(mockRepo, mockGen)
			got, err := svc.ShortlengthURL(gc, tc.originURL, tc.expireAt)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, got)
			}

			// Thêm: AssertExpectations
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
		setupMock   func(mockRepo *mockrepo.UrlStorage, ctx *gin.Context)
		expectedRes string
		expectedErr bool
	}{
		{
			name: "normal case",
			code: "abc1234",
			setupMock: func(mockRepo *mockrepo.UrlStorage, ctx *gin.Context) {
				mockRepo.On("GetUrl", ctx, "abc1234").Return("https://example.com", nil)
			},
			expectedRes: "https://example.com",
			expectedErr: false,
		},
		{
			name: "url not found",
			code: "invalid",
			setupMock: func(mockRepo *mockrepo.UrlStorage, ctx *gin.Context) {
				mockRepo.On("GetUrl", ctx, "invalid").Return("", errors.New("url not found"))
			},
			expectedRes: "",
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gin.SetMode(gin.TestMode)

			mockRepo := mockrepo.NewUrlStorage(t)
			mockGen := mockcodegen.NewCodeGenerator(t)
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupMock(mockRepo, gc)

			svc := NewShortenURL(mockRepo, mockGen)
			got, err := svc.GetURL(gc, tc.code)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
