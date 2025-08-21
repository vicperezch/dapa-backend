package model

import (
	"time"
	"gorm.io/gorm"
)

type FormStatus string

const (
	FormStatusPending   FormStatus = "pending"
	FormStatusCancelled FormStatus = "cancelled"
	FormStatusApproved  FormStatus = "approved"
)

type User struct {
	ID                    uint       `json:"id" gorm:"primaryKey"`
	Name                  string     `json:"name" gorm:"size:50;not null"`
	LastName              string     `json:"lastName" gorm:"column:last_name;size:50;not null"`
	Phone                 string     `json:"phone" gorm:"size:20;not null"`
	Email                 string     `json:"email" gorm:"default:unknown;unique"`
	PasswordHash          string     `json:"password" gorm:"column:password_hash"`
	Role                  string     `json:"role" gorm:"size:20;not null"`
	LicenseExpirationDate time.Time  `json:"licenseExpirationDate" gorm:"column:license_expiration_date"`
	CreatedAt             time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	LastModifiedAt        time.Time  `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt             *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive              bool       `json:"isActive" gorm:"column:is_active;default:true"`
}

type Vehicle struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	Brand          string     `json:"brand" gorm:"size:50;not null"`
	Model          string     `json:"model" gorm:"size:50;not null"`
	LicensePlate   string     `json:"licensePlate" gorm:"size:15;not null;unique"`
	CapacityKg     float64    `json:"capacityKg" gorm:"column:capacity_kg"`
	Available      bool       `json:"available"`
	InsuranceDate  time.Time  `json:"insuranceDate" gorm:"column:insurance"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	LastModifiedAt time.Time  `json:"lastModifiedAt" gorm:"column:last_modified_at;autoUpdateTime"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	IsActive       bool       `json:"isActive" gorm:"column:is_active;default:true"`
}

type Quote struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	SubmissionID uint    `gorm:"column:submission_id"`
	DriverID     uint    `gorm:"column:driver_id"`
	VehicleID    uint    `gorm:"column:vehicle_id"`
	StateID      uint    `gorm:"column:state_id"`
	ServiceType  uint    `gorm:"column:service_type"`
	TotalAmount  float64 `gorm:"column:total_amount"`
	Date         time.Time
	Details      string
}

type QuoteView struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Date         time.Time `json:"date"`
	Client       string    `json:"clientName"`
	VehicleBrand string    `json:"vehicleBrand" gorm:"column:vehicle_brand"`
	VehicleModel string    `json:"vehicleModel" gorm:"column:vehicle_model"`
	Driver       string    `json:"driver"`
}

func (QuoteView) TableName() string {
	return "pending_orders"
}

// ******************** FORMULARIO ********************
// Tipos de preguntas disponibles
type QuestionType struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Type string `json:"type" gorm:"size:50;not null" validate:"required,question_type"`
}

// Pregunta del formulario
type Question struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Question    string           `json:"question" gorm:"size:50;not null" validate:"required,question_text"`
	Description *string          `json:"description,omitempty" gorm:"size:255" validate:"omitempty,question_desc"`
	TypeID      uint             `json:"typeId" gorm:"not null" validate:"required,gt=0"`
	IsActive    bool             `json:"isActive" gorm:"not null;default:true"`
	Position    int              `json:"position" gorm:"not null;default:1"`
	Type        QuestionType     `json:"type" gorm:"foreignKey:TypeID"`
	Options     []QuestionOption `json:"options,omitempty" gorm:"foreignKey:QuestionID"`
	DeletedAt   gorm.DeletedAt   `json:"deletedAt,omitempty" gorm:"index"`
}


// Opciones de una pregunta
type QuestionOption struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	QuestionID uint   `json:"questionId" gorm:"not null" validate:"required,gt=0"`
	Option     string `json:"option" gorm:"size:50;not null" validate:"required,question_option"`
	DeletedAt  gorm.DeletedAt   `json:"deletedAt,omitempty" gorm:"index"`
}

// Envío de formulario
type Submission struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	SubmittedAt time.Time  `json:"submittedAt" gorm:"default:CURRENT_TIMESTAMP"`
	Status      FormStatus `json:"status" gorm:"type:form_status;not null;default:'pending'" validate:"required,submission_status"`
	Answers     []Answer   `json:"answers,omitempty" gorm:"foreignKey:SubmissionID"`
}

// Respuesta a una pregunta
type Answer struct {
	ID           uint             `json:"id" gorm:"primaryKey"`
	SubmissionID uint             `json:"submissionId" gorm:"not null" validate:"required,gt=0"`
	QuestionID   uint             `json:"questionId,omitempty" validate:"omitempty,gt=0"`
	Answer       *string          `json:"answer,omitempty" validate:"omitempty,max=255"`
	Question     Question         `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	Options      []QuestionOption `json:"options,omitempty" gorm:"many2many:answer_options;"`
}

type SubmissionStats struct {
    TotalSubmissions     int64                       `json:"totalSubmissions"`
    SubmissionsByStatus  []StatusCount               `json:"submissionsByStatus"`
    AnswersByQuestion    []QuestionAnswerDistribution `json:"answersByQuestion"`
}

type StatusCount struct {
    Status string `json:"status"`
    Count  int64  `json:"count"`
}

type QuestionAnswerDistribution struct {
    QuestionID uint  `json:"questionId"`
    OptionID   uint  `json:"optionId"`
    Count      int64 `json:"count"`
}

