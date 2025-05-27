package utils

import "golang.org/x/crypto/bcrypt"

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