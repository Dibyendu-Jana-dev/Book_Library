package utility

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

func IsStrongPassword(password string) bool {
	// Check for at least one special character
	specialCharRegex := `[[:punct:]]`
	if ok, _ := regexp.MatchString(specialCharRegex, password); !ok {
		return false
	}
	// Check for at least one capital letter
	upperCaseRegex := `[[:upper:]]`
	if ok, _ := regexp.MatchString(upperCaseRegex, password); !ok {
		return false
	}
	// Check for at least one digit
	digitRegex := `[[:digit:]]`
	if ok, _ := regexp.MatchString(digitRegex, password); !ok {
		return false
	}
	// Check for minimum length of 8
	if len(password) < 8 {
		return false
	}

	return true
}

func HashPassword(password string) (string, error) {
	// Hash the password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(hashedPassword, enteredPassword string) error {
	// Compare the hashed password with the entered password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(enteredPassword))
	return err
}