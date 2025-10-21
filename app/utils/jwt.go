package utils

import (
	"time"

	"dapa/app/model"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(EnvMustGet("JWT_SECRET"))

// Genera un token JWT para un usuario
// Recibe el usuario como parámetro
// Retorna el token como string
func GenerateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &model.EmployeeClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Determina si un token JWT es válido
// Retorna las claims si es válido
// Retorna un error si no lo es
func ValidateToken(tokenString string) (*model.EmployeeClaims, error) {
	claims := &model.EmployeeClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
