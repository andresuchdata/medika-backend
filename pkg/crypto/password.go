package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default bcrypt cost
	DefaultCost = bcrypt.DefaultCost
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// MustHashPassword hashes a password and panics on error (for testing/seeding only)
func MustHashPassword(password string) string {
	hash, err := HashPassword(password)
	if err != nil {
		panic(fmt.Sprintf("failed to hash password: %v", err))
	}
	return hash
}
