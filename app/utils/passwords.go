package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

var resetSecret = []byte(EnvMustGet("RESET_SECRET"))

// Hashea una cadena en texto plano utilizando bcrypt
// Retorna el hash como string o un error
func HashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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
