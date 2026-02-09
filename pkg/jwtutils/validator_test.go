package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTValidator_ValidateToken(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		publicKeyPath  string
		tokenString    string
		expectedClaims jwt.MapClaims
		expectedErr    error
	}{
		{
			name:          "normal case",
			publicKeyPath: filepath.FromSlash("./public.test.pem"),
			tokenString:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQiLCJuYW1lIjoiSm9obiJ9.Tj6o3td461ey04A_CdovTMohilao1N9cKEpNOH9pbiHw5CFTSOUxnKttzYv6P0Kgm-8Q0cdYRgt0M_LqVaGhprt-Bjte8iskkSN1M336o22KZmg0K4k3x0n1pccHDF6sboJv5-krjZEG_XxKazHAOQCt02TlHPK9XEdYJJRD-iyQVHDw4el4lRwteR1bDtZV6HqBTgUf6EHgbSQWMsoVtxmNy3ex8dzySPKU-F2A6Y8TiM9DzVi0PZFvV5I8Jn_0CKB2c1WiLcr2AbjAfct8KyqAwdub462mxeeUjV8EiOrfMZKnA0E9imu83GVg9EhnHvg2pfdJALBIUJmaYpK-Lw",
			expectedClaims: jwt.MapClaims{
				"id":   "1234",
				"name": "John",
			},
			expectedErr: nil,
		},
		{
			name:           "invalid token string",
			publicKeyPath:  filepath.FromSlash("./public.test.pem"),
			tokenString:    "this.is.not.a.jwt",
			expectedClaims: nil,
			expectedErr:    errInvalidToken,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator, err := NewJWTValidator(tc.publicKeyPath)
			assert.NoError(t, err)

			claims, err := validator.ValidateToken(tc.tokenString)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedClaims, claims)
		})
	}
}
