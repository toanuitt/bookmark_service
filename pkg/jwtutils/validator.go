package jwtutils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name JWTValidator --filename jwt_validator.go
var errInvalidToken = errors.New("invalid token")

type JWTValidator interface {
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type jwtValidator struct {
	publicKey *rsa.PublicKey
}

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

func (j *jwtValidator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}

	return token.Claims.(jwt.MapClaims), nil
}
