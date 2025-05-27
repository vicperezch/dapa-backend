package utils

import (
	"time"

	"dapa/app/model"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret holds the secret key used for signing JWT tokens.
// It is loaded from the environment variable JWT_SECRET.
var jwtSecret = []byte(EnvMustGet("JWT_SECRET"))

// GenerateToken creates a JWT token for the given employee.
// The token includes employee details as claims and expires after 24 hours.
func GenerateToken(employee *model.Employee) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &model.EmployeeClaims{
		EmployeeID: employee.ID,
		UserID:     employee.UserID,
		Name:       employee.User.Name + " " + employee.User.LastName,
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

// ValidateToken parses and validates a JWT token string.
// It returns the EmployeeClaims if the token is valid, or an error otherwise.
func ValidateToken(tokenString string) (*model.EmployeeClaims, error) {
	claims := &model.EmployeeClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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