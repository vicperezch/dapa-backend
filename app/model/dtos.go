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

type CreateVehicleRequest struct {
	Brand                  string  `json:"brand" binding:"required"`
	Model                  string  `json:"model" binding:"required"`
	LicensePlate           string  `json:"licensePlate" binding:"required"`
	CapacityKg             float64 `json:"capacityKg" binding:"required,gt=0"`
	Available              *bool   `json:"available" binding:"required"` 
	CurrentMileage         float64 `json:"currentMileage" binding:"required,gt=0"`
	NextMaintenanceMileage float64 `json:"nextMaintenanceMileage" binding:"required,gt=0"`
}

type UpdateVehicleRequest struct {
	Brand                  string  `json:"brand" binding:"required"`
	Model                  string  `json:"model" binding:"required"`
	LicensePlate           string  `json:"licensePlate" binding:"required"` 
	CapacityKg             float64 `json:"capacityKg" binding:"required,gt=0"`
	Available              *bool   `json:"available" binding:"required"` 
	CurrentMileage         float64 `json:"currentMileage" binding:"required,gt=0"`
	NextMaintenanceMileage float64 `json:"nextMaintenanceMileage" binding:"required,gt=0"`
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
