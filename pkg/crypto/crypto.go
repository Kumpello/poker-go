package crypto

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("cannot encrypt the password: %w", err)
	}

	return string(bytes), nil
}

//VerifyPassword compares user password with the encrypted one
func VerifyPassword(userPassword string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))

	if err != nil {
		return fmt.Errorf("password is incorrect: %w", err)
	}

	return nil
}
