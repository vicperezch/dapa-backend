package model

import "time"

type UserDTO struct {
	Name                  string    `json:"name" binding:"required"`
	LastName              string    `json:"lastName" binding:"required"`
	Phone                 string    `json:"phone" binding:"required,phone"`
	Email                 string    `json:"email" binding:"required,email"`
	Role                  string    `json:"role" binding:"required,oneof=admin driver"`
	LicenseExpirationDate time.Time `json:"licenseExpirationDate"`
}

type VehicleDTO struct {
	Brand         string    `json:"brand" binding:"required"`
	Model         string    `json:"model" binding:"required"`
	LicensePlate  string    `json:"licensePlate" binding:"required,plate"`
	CapacityKg    float64   `json:"capacityKg" binding:"required,gt=0"`
	IsAvailable   bool      `json:"isAvailable" binding:"required"`
	InsuranceDate time.Time `json:"insuranceDate" binding:"required"`
}

type OrderDTO struct {
	UserID      *uint   `json:"userId"`
	VehicleID   *uint   `json:"vehicleId"`
	ClientName  string  `json:"clientName" binding:"required"`
	ClientPhone string  `json:"clientPhone" binding:"required"`
	Origin      string  `json:"origin" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	TotalAmount float64 `json:"totalAmount" binding:"required"`
	Details     *string `json:"details"`
	Type        string  `json:"type" binding:"required"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type RegisterDTO struct {
	Name                  string    `json:"name" binding:"required"`
	LastName              string    `json:"lastName" binding:"required"`
	Phone                 string    `json:"phone" binding:"required,phone"`
	Email                 string    `json:"email" binding:"required,email"`
	LicenseExpirationDate time.Time `json:"licenseExpirationDate" binding:"required_if=Role driver"`
	Password              string    `json:"password" binding:"required,password"`
	Role                  string    `json:"role" binding:"required,oneof=admin driver"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordDTO struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,password"`
}

type AcceptSubmissionDTO struct {
	SubmissionID uint    `json:"submissionId" binding:"required"`
	ClientName   string  `json:"clientName" binding:"required"`
	ClientPhone  string  `json:"clientPhone" binding:"required"`
	Origin       string  `json:"origin" binding:"required"`
	Destination  string  `json:"destination" binding:"required"`
	TotalAmount  float64 `json:"totalAmount" binding:"required"`
	Details      string  `json:"details" binding:"required"`
	Type         string  `json:"type" binding:"required"`
}

type AssignOrderDTO struct {
	UserID    uint `json:"userId" binding:"required"`
	VehicleID uint `json:"vehicleId" binding:"required"`
}

type OrderStatusDTO struct {
	Status string `json:"status" binding:"oneof=pending assigned pickup collected delivered"`
}

type QuestionTypeDTO struct {
	Type string `json:"type" binding:"required,max=50"`
}

type QuestionOptionDTO struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	QuestionID uint   `json:"questionId" gorm:"not null" validate:"required,gt=0"`
	Option     string `json:"option" gorm:"size:50;not null" validate:"required,question_option"`
}

type QuestionDTO struct {
	Question    string              `json:"question" binding:"required,max=50"`
	Description *string             `json:"description,omitempty" binding:"omitempty,max=255"`
	TypeID      uint                `json:"typeId" binding:"required"`
	IsActive    *bool               `json:"isActive,omitempty"`
	Options     []QuestionOptionDTO `json:"options,omitempty"`
}

type ReorderQuestionDTO struct {
	SourceID uint `json:"sourceId" binding:"required"`
	TargetID uint `json:"targetId" binding:"required"`
}

type AnswerDTO struct {
	QuestionID uint    `json:"questionId,omitempty"`
	Answer     *string `json:"answer,omitempty"`
	OptionsID  []uint  `json:"optionsId,omitempty"`
}

type CreateSubmissionDTO struct {
	Answers []AnswerDTO `json:"answers" binding:"required,min=1"`
}

type UpdateSubmissionStatusDTO struct {
	Status FormStatus `json:"status" binding:"required,oneof=pending cancelled approved"`
}

type OrderTokenDTO struct {
	Token string `json:"token" binding:"required"`
}

type OrderTrackingDTO struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}

type FinancialReportDTO struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	TotalAmount float64   `json:"totalAmount"`
	User        string    `json:"user"`
}

type DriverReportDTO struct {
	DriverName    string  `json:"driverName"`
	TotalOrders   int     `json:"totalOrders"`
	OrdersPerWeek float64 `json:"ordersPerWeek"`
}

type TotalIncomeReportDTO struct {
	TotalIncome float64 `json:"totalIncome"`
}

type CompletedQuotationsDTO struct {
	Series     []struct {
		Data []int  `json:"data"`
	} `json:"series"`
	Categories []string `json:"categories"`
}

type QuotationsStatusDTO struct {
	Series []float64 `json:"series"`
	Labels []string  `json:"labels"`
}

type DriverPerformanceDTO struct {
	Series     []struct {
		Data []float64 `json:"data"`
	} `json:"series"`
	Categories []string `json:"categories"`
}

type DriversBarDataDTO struct {
	Series []struct {
		Name string `json:"name"`
		Data []int  `json:"data"`
	} `json:"series"`
	Categories []string `json:"categories"`
}

type TripParticipationDTO struct {
	Series []int    `json:"series"`
	Labels []string `json:"labels"`
}

type ApiResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    any      `json:"data"`
	Errors  []string `json:"errors"`
}
