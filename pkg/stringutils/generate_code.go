package stringutils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

//go:generate mockery --name CodeGenerator --filename generate_code.go
type CodeGenerator interface {
	GenerateCode(length int) (string, error)
}

type codeGeneratorPass struct {
}

func NewCodeGenerator() CodeGenerator {
	return &codeGeneratorPass{}
}

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
