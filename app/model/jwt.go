package model

import "github.com/golang-jwt/jwt/v5"

type EmployeeClaims struct {
	EmployeeID uint   `json:"id"`
	UserID     uint   
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}