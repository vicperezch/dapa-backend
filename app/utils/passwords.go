package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

var resetSecret = []byte(EnvMustGet("RESET_SECRET"))

// Hashea una contraseña en texto plano utilizando bcrypt
// Retorna el hash como string o un error
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Compara una contraseña en texto plano contra un hash
// Retorna el valor de la comparación como boolean
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Genera un token válido para reestablecer contraseña
// Retorna el token como string o un error
func GenerateResetToken(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
