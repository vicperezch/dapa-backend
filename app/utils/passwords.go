package utils

import "golang.org/x/crypto/bcrypt"

// Utiliza bcrypt para hashear una contraseña
// Retorna el hash en forma de string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

// Valida si una contraseña en texto plano y un hash son equivalentes
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
