package jwtutils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name JWTValidator --filename jwt_validator.go

// errInvalidToken is returned when a token is invalid, expired, or cannot be verified.
var errInvalidToken = errors.New("invalid token")

// JWTValidator defines an interface for validating JWT tokens and extracting claims.
// Implementations are responsible for verifying the token signature and integrity.
type JWTValidator interface {
	// ValidateToken verifies the given token string and returns its claims if valid.
	// It returns an error if the token is invalid, expired, or cannot be parsed.
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

// jwtValidator is the concrete implementation of JWTValidator.
// It validates JWT tokens using an RSA public key.
type jwtValidator struct {
	publicKey *rsa.PublicKey
}

// NewJWTValidator creates a new JWTValidator using the RSA public key located at publicKeyPath.
//
// The public key file must be in PEM format. An error is returned if the file
// cannot be read or the key cannot be parsed.
func NewJWTValidator(publicKeyPath string) (JWTValidator, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtValidator{publicKey: publicKey}, nil
}

// ValidateToken parses and validates the given JWT token string using the configured
// RSA public key.
//
// If the token is valid, it returns the token claims as jwt.MapClaims.
// If the token is invalid, expired, or cannot be verified, it returns errInvalidToken.
func (j *jwtValidator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}

	return token.Claims.(jwt.MapClaims), nil
}
