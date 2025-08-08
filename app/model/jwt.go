package model

import "github.com/golang-jwt/jwt/v5"

type EmployeeClaims struct {
	UserID uint
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

