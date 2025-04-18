package utils

import (
	"time"

	"dapa/app/model"

	"github.com/golang-jwt/jwt/v5"
)

//Se obtiene la clave secreta para JWT del archivo .env
var jwtSecret = []byte(EnvMustGet("JWT_SECRET"))

//Genera un token para un empleado
func GenerateToken(employee *model.Employee) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &model.EmployeeClaims{
		EmployeeID: employee.ID,
		UserID:     employee.UserID,
		Name:   	employee.User.Name + " " + employee.User.LastName,
		Email:      employee.User.Email,
		Role:       employee.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

//Valida un token y obtiene los claims a partir de este
func ValidateToken(tokenString string) (*model.EmployeeClaims, error) {
	claims := &model.EmployeeClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}



