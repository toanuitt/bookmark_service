package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTGenerator_GenerateToken(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name           string
		keyPath        string
		inputContent   jwt.MapClaims
		expectedOutput string
		expectedErr    error
	}{
		{
			name:    "normal case",
			keyPath: filepath.FromSlash("./private.test.pem"),
			inputContent: jwt.MapClaims{
				"id":   "1234",
				"name": "John",
			},
			expectedErr:    nil,
			expectedOutput: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQiLCJuYW1lIjoiSm9obiJ9.Tj6o3td461ey04A_CdovTMohilao1N9cKEpNOH9pbiHw5CFTSOUxnKttzYv6P0Kgm-8Q0cdYRgt0M_LqVaGhprt-Bjte8iskkSN1M336o22KZmg0K4k3x0n1pccHDF6sboJv5-krjZEG_XxKazHAOQCt02TlHPK9XEdYJJRD-iyQVHDw4el4lRwteR1bDtZV6HqBTgUf6EHgbSQWMsoVtxmNy3ex8dzySPKU-F2A6Y8TiM9DzVi0PZFvV5I8Jn_0CKB2c1WiLcr2AbjAfct8KyqAwdub462mxeeUjV8EiOrfMZKnA0E9imu83GVg9EhnHvg2pfdJALBIUJmaYpK-Lw",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testGen, err := NewJWTGenerator(tc.keyPath)
			assert.Equal(t, err, tc.expectedErr)
			res, err := testGen.GenerateToken(tc.inputContent)
			assert.Equal(t, res, tc.expectedOutput)
			assert.Equal(t, err, tc.expectedErr)
		})
	}
}
