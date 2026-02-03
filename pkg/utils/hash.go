package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain-text password using bcrypt with the default cost.
// It returns the hashed password as a string.
//
// Note:
// - This function ignores the error from bcrypt.GenerateFromPassword for simplicity.
// - In production code, you may want to return (string, error) instead to handle failures properly.
func HashPassword(password string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes)
}

// VerifyPassword compares a plain-text password with a bcrypt hashed password.
// It returns true if the password matches the hash, otherwise false.
func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
