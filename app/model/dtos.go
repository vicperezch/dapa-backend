package model

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"lastName" binding:"required"`
	Phone    string `json:"phone" binding:"required,phone"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
	Role     string `json:"role" binding:"required,validrole"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"lastName" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"lastName" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email"`
	Role     string `json:"role" binding:"required,validrole"`
}

type ApiResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}