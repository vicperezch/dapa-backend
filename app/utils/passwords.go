package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

var resetSecret = []byte(EnvMustGet("RESET_SECRET"))

// Hashea una contaseña en texto plano utilizando bcrypt
// Retorna el hash como string o un error
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Hashea una cadena en texto plano utilizando bcrypt
// Retorna el hash como string
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Compara una contraseña en texto plano contra un hash
// Retorna el valor de la comparación como boolean
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Genera un token válido y seguro
// Retorna el token como string o un error
func GenerateSecureToken(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
