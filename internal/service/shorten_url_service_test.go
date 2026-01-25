package service

import (
	"context"
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
		setupMock   func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context)
		expectedRes string
		expectedErr bool
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
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
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
		})
	}
}

func TestShortenURL_GetURL(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		code        string
		setupMock   func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context)
		expectedRes string
		expedtedErr bool
	}{
		{
			name: "normal case",
			code: "abc1234",
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context) {
				mockRepo.On("GetUrl", ctx, "abc1234").Return("https://example.com", nil)
			},
			expectedRes: "https://example.com",
			expedtedErr: false,
		},
		{
			name: "url not found",
			code: "invalid",
			setupMock: func(mockRepo *mockrepo.UrlStorage, mockGen *mockcodegen.CodeGenerator, ctx context.Context) {
				mockRepo.On("GetUrl", ctx, "invalid").Return("", errors.New("url not found"))
			},
			expectedRes: "",
			expedtedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mockrepo.NewUrlStorage(t)
			mockGen := mockcodegen.NewCodeGenerator(t)
			ctx := context.Background()
			tc.setupMock(mockRepo, mockGen, ctx)

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
