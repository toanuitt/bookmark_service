package jwtutils

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name JWTGenerator --filename jwt_generator.go

// JWTGenerator defines an interface for generating signed JWT tokens.
// Implementations are responsible for signing tokens using a private key.
type JWTGenerator interface {
	// GenerateToken creates and signs a JWT token with the given claims.
	// It returns the signed token string or an error if signing fails.
	GenerateToken(jwtContent jwt.MapClaims) (string, error)
}

// jwtGenerator is the concrete implementation of JWTGenerator.
// It signs JWT tokens using an RSA private key.
type jwtGenerator struct {
	privateKey *rsa.PrivateKey
}

// NewJWTGenerator creates a new JWTGenerator using the RSA private key located at privateKeyPath.
//
// The private key file must be in PEM format. An error is returned if the file
// cannot be read or the key cannot be parsed.
func NewJWTGenerator(privateKeyPath string) (JWTGenerator, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtGenerator{privateKey: privateKey}, nil
}

// GenerateToken creates a new JWT using the provided claims and signs it
// with the RSA private key configured in the generator.
//
// It returns the signed token string or an error if signing fails.
func (g *jwtGenerator) GenerateToken(jwtContent jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtContent)
	tokenString, err := token.SignedString(g.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
