package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password
func HashPassword(pw string) (string, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPw), nil
}

// CheckPassword check if the provided pw is correct.
func CheckPassword(pw, hashedPw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw))
}
