package model

import "github.com/golang-jwt/jwt/v5"

type EmployeeClaims struct {
	UserID uint
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
