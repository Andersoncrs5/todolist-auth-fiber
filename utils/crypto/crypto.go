package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Encoder(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("Password is Required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", fmt.Errorf("Error the encoder password")
	}

	return string(hash), nil
}

func Compare(password string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false
	}

	return true
}