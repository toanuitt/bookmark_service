package service

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	passLength = 10
)

// passwordService is the concrete implementation of the Password interface.
// It uses cryptographically secure random number generation to create passwords.
type passwordService struct {
}

//go:generate mockery --name Password --filename pass_service.go

// Password defines the interface for password generation operations.
type Password interface {
	GeneratePassword() (string, error)
}

// NewPassword creates and returns a new Password service instance.
func NewPassword() Password {
	return &passwordService{}
}

// GeneratePassword generates a cryptographically secure random password.
// It uses crypto/rand to select random characters from the charset with a fixed length of 10 characters.
// Returns an error if random number generation fails.
func (s *passwordService) GeneratePassword() (string, error) {
	var strBuilder bytes.Buffer

	for range passLength {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}
	return strBuilder.String(), nil
}
