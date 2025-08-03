package utils

import (
	"time"

	"github.com/dchest/passwordreset"
	"golang.org/x/crypto/bcrypt"
)

var resetSecret = []byte(EnvMustGet("RESET_SECRET"))

// HashPassword hashes a plain text password using bcrypt algorithm.
// Returns the hashed password as a string or an error if hashing fails.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword compares a plain text password with a hashed password.
// Returns true if they match, false otherwise.
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generates a new token to reset a users password
// Returns the generated token
func GenerateResetToken(email string, passwordHash string) string {
	token := passwordreset.NewToken(email, time.Hour, []byte(passwordHash), resetSecret)

	return token
}

// Checks if a token to reset a users password is valid
// Returns the user that requested the reset
func VerifyResetToken(token string, getPasswordHash func(string) ([]byte, error)) (string, error) {
	login, err := passwordreset.VerifyToken(token, getPasswordHash, resetSecret)

	if err != nil {
		return "", err
	}

	return login, nil
}

