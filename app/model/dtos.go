package model

import "time"

type RegisterRequest struct {
	Name              string    `json:"name" binding:"required"`
	LastName          string    `json:"lastName" binding:"required"`
	Phone             string    `json:"phone" binding:"required,phone"`
	Email             string    `json:"email" binding:"required,email"`
	LicenseExpiration time.Time `json:"licenseExpirationDate"`
	Password          string    `json:"password" binding:"required,password"`
	Role              string    `json:"role" binding:"required,validrole"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type NewPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,password"`
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

type CreateQuestionTypeRequest struct {
	Type string `json:"type" binding:"required,max=50"`
}

type UpdateQuestionTypeRequest struct {
	Type string `json:"type" binding:"required,max=50"`
}

type QuestionOptionRequest struct {
	Option string `json:"option" binding:"required,max=50"`
}

type CreateQuestionRequest struct {
	Question    string                  `json:"question" binding:"required,max=50"`
	Description *string                 `json:"description,omitempty" binding:"omitempty,max=255"`
	TypeID      uint                    `json:"typeId" binding:"required"`
	IsActive    *bool                   `json:"isActive,omitempty"`
	Options     []QuestionOptionRequest `json:"options,omitempty"`
}

type UpdateQuestionRequest struct {
	Question    *string                 `json:"question,omitempty" binding:"omitempty,max=50"`
	Description *string                 `json:"description,omitempty" binding:"omitempty,max=255"`
	TypeID      *uint                   `json:"typeId,omitempty"`
	IsActive    *bool                   `json:"isActive,omitempty"`
	Options     []QuestionOptionRequest `json:"options,omitempty"`
}

type AnswerRequest struct {
	QuestionID *uint   `json:"questionId,omitempty"`
	Answer     *string `json:"answer,omitempty"`
	OptionID   *uint   `json:"optionId,omitempty"`
}

type CreateSubmissionRequest struct {
	UserID  uint            `json:"userId" binding:"required"`
	Answers []AnswerRequest `json:"answers" binding:"required,min=1"`
}

type UpdateSubmissionStatusRequest struct {
	Status FormStatus `json:"status" binding:"required,oneof=pending cancelled approved"`
}

type QuestionResponse struct {
	ID          uint                   `json:"id"`
	Question    string                 `json:"question"`
	Description *string                `json:"description,omitempty"`
	TypeID      uint                   `json:"typeId"`
	Type        string                 `json:"type"`
	IsActive    bool                   `json:"isActive"`
	Options     []QuestionOptionResponse `json:"options,omitempty"`
}

type QuestionOptionResponse struct {
	ID     uint   `json:"id"`
	Option string `json:"option"`
}

type SubmissionResponse struct {
	ID          uint             `json:"id"`
	UserID      uint             `json:"userId"`
	UserName    string           `json:"userName"`
	UserEmail   string           `json:"userEmail"`
	SubmittedAt string           `json:"submittedAt"`
	Status      FormStatus       `json:"status"`
	Answers     []AnswerResponse `json:"answers,omitempty"`
}

type AnswerResponse struct {
	ID           uint    `json:"id"`
	QuestionID   *uint   `json:"questionId,omitempty"`
	Question     *string `json:"question,omitempty"`
	Answer       *string `json:"answer,omitempty"`
	OptionID     *uint   `json:"optionId,omitempty"`
	OptionText   *string `json:"optionText,omitempty"`
}

type QuestionFilters struct {
	TypeID   *uint  `form:"typeId"`
	IsActive *bool  `form:"isActive"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	Limit    int    `form:"limit,default=10" binding:"min=1,max=100"`
}

type SubmissionFilters struct {
	UserID *uint      `form:"userId"`
	Status *FormStatus `form:"status" binding:"omitempty,oneof=pending cancelled approved"`
	Page   int        `form:"page,default=1" binding:"min=1"`
	Limit  int        `form:"limit,default=10" binding:"min=1,max=100"`
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
