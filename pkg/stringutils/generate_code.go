package stringutils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// CodeGenerator defines the interface for generating random codes.
//
//go:generate mockery --name CodeGenerator --filename generate_code.go
type CodeGenerator interface {
	GenerateCode(length int) (string, error)
}

// codeGeneratorPass is the concrete implementation of the CodeGenerator interface.
type codeGeneratorPass struct {
}

// NewCodeGenerator creates and returns a new CodeGenerator instance.
func NewCodeGenerator() CodeGenerator {
	return &codeGeneratorPass{}
}

// GenerateCode generates a random alphanumeric code of the specified length.
// It uses cryptographically secure random number generation.
func (s *codeGeneratorPass) GenerateCode(length int) (string, error) {
	var strBuilder bytes.Buffer

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}
	return strBuilder.String(), nil
}
